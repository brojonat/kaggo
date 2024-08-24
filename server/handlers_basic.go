package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func writeOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(api.DefaultJSONResponse{Message: "ok"})
}

func writeInternalError(l *slog.Logger, w http.ResponseWriter, e error) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError, e.Error(), pcs[0])
	_ = l.Handler().Handle(context.Background(), r)
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(api.DefaultJSONResponse{Error: "internal error"})
}

func writeBadRequestError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	resp := api.DefaultJSONResponse{Error: err.Error()}
	json.NewEncoder(w).Encode(resp)
}

func writeEmptyResultError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	resp := api.DefaultJSONResponse{Error: "empty result set"}
	json.NewEncoder(w).Encode(resp)
}

// handlePing pings the database
func handlePing(l *slog.Logger, p *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := p.Ping(r.Context())
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		json.NewEncoder(w).Encode(api.DefaultJSONResponse{Message: "PONG"})
	}
}

// handleGetToken returns a token
func handleIssueToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, pwd, ok := r.BasicAuth()
		if !ok {
			writeBadRequestError(w, fmt.Errorf("bad credentials"))
			return
		}
		if pwd != getSecretKey() {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.DefaultJSONResponse{Error: "not authorized"})
			return
		}
		sc := jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * 7 * 24 * time.Hour).Unix(),
		}
		c := authJWTClaims{
			StandardClaims: sc,
			Email:          email,
		}
		token, _ := generateAccessToken(c)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(api.DefaultJSONResponse{Message: token})
	}
}
