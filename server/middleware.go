package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/server-tools/stools"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
)

const FirebaseJWTHeader = "Firebase-JWT"

type contextKey int

var jwtCtxKey contextKey = 1

// Convenience middleware that applies commonly used middleware to the wrapped
// handler. This will make the handler gracefully handle panics, sets the
// content type to application/json, limits the body size that clients can send,
// wraps the handler with the usual CORS settings.
func apiMode(l *slog.Logger, maxBytes int64, headers, methods, origins []string) stools.HandlerAdapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next = makeGraceful(l)(next)
			next = setMaxBytesReader(maxBytes)(next)
			next = setContentType("application/json")(next)
			handlers.CORS(
				handlers.AllowedHeaders(headers),
				handlers.AllowedMethods(methods),
				handlers.AllowedOrigins(origins),
			)(next).ServeHTTP(w, r)
		}
	}
}

func setContentType(content string) stools.HandlerAdapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", content)
			next(w, r)
		}
	}
}

func makeGraceful(l *slog.Logger) stools.HandlerAdapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					l.Error("recovered from panic")
					switch v := err.(type) {
					case error:
						writeInternalError(l, w, v)
					case string:
						writeInternalError(l, w, fmt.Errorf(v))
					default:
						writeInternalError(l, w, fmt.Errorf("recovered but unexpected type from recover()"))
					}
				}
			}()
			next.ServeHTTP(w, r)
		}
	}
}

func setMaxBytesReader(mb int64) stools.HandlerAdapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, mb)
			next(w, r)
		}
	}
}

func bearerAuthorizer(gsk func() string) func(*http.Request) bool {
	return func(r *http.Request) bool {
		var claims authJWTClaims
		ts := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if ts == "" {
			return false
		}
		kf := func(token *jwt.Token) (interface{}, error) {
			return []byte(gsk()), nil
		}
		token, err := jwt.ParseWithClaims(ts, &claims, kf)
		if err != nil || !token.Valid {
			return false
		}
		ctx := context.WithValue(r.Context(), jwtCtxKey, token.Claims)
		*r = *r.WithContext(ctx)
		return true
	}
}

// Iterates over the supplied authorizers and if at least one passes, then the
// next handler is called, otherwise an unauthorized response is written.
func atLeastOneAuth(authorizers ...func(*http.Request) bool) stools.HandlerAdapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for _, a := range authorizers {
				if !a(r) {
					continue
				}
				next(w, r)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.DefaultJSONResponse{Error: "unauthorized"})
		}
	}
}

// Increment a prometheus counter for each request
func withPromCounter(pm *prometheus.CounterVec) stools.HandlerAdapter {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			nwr := negroni.NewResponseWriter(w)
			hf(nwr, r)
			name := fmt.Sprintf("%s %s", r.Method, r.URL)
			sc := fmt.Sprintf("%d", nwr.Status())
			pm.With(prometheus.Labels{"name": name, "code": sc}).Inc()
		}
	}
}
