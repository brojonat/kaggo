package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/brojonat/kaggo/server/db/dbgen"
	"github.com/brojonat/server-tools/stools"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.temporal.io/sdk/client"
)

// These are prometheus metric keys that root handler expects to be passed it.
// When the root handler is constructed, it passes these metrics in to
// particular handlers. If the caller misconfigures the prometheus metrics
// passed in to the root handler, it will panic (on startup, so it's fine).
const (
	MetricKeyInternalRandom = "internal-random"
)

// This is a convenience method for getting the necessary metrics. Some handlers
// assume that particular metrics are passed in as particular types, but any
// programming error on behalf of the caller will panic this at startup, so that
// should be tolerable for now.
func GetDefaultPromMetrics() map[string]prometheus.Collector {
	return map[string]prometheus.Collector{
		MetricKeyInternalRandom: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "internal_random",
				Help: "A pseudo random metric",
			},
			[]string{"id"},
		),
	}
}

// Run the HTTP server. Note that promMetrics in general shouldn't be configured by the
// caller. The caller should simply use the value returned by GetDefaultPromMetrics.
// In some specialized cases or in tests, it may be useful to pass in custom values,
// but doing so means that you MUST be prepared to handle this
func RunHTTPServer(
	ctx context.Context,
	port string,
	l *slog.Logger,
	dbHost string,
	tcHost string,
	promMetrics map[string]prometheus.Collector,
) error {

	p, err := stools.GetConnPool(
		ctx, dbHost,
		func(ctx context.Context, c *pgx.Conn) error { return nil },
	)
	if err != nil {
		return fmt.Errorf("could not connect to db: %w", err)
	}
	q := dbgen.New(p)

	tc, err := client.Dial(client.Options{
		Logger:   l,
		HostPort: os.Getenv("TEMPORAL_HOST"),
	})
	if err != nil {
		return fmt.Errorf("could not initialize Temporal client: %w", err)
	}
	defer tc.Close()

	pms := []prometheus.Collector{}
	for _, pm := range promMetrics {
		pms = append(pms, pm)
	}
	prometheus.MustRegister(pms...)
	router, err := getRouter(l, p, q, tc, promMetrics)
	if err != nil {
		return err
	}

	listenAddr := fmt.Sprintf(":%s", port)
	l.Info("listening", "port", listenAddr)
	return http.ListenAndServe(listenAddr, router)
}

func getRouter(
	l *slog.Logger,
	p *pgxpool.Pool,
	q *dbgen.Queries,
	tc client.Client,
	promMetrics map[string]prometheus.Collector,
) (http.Handler, error) {
	// new router
	mux := http.NewServeMux()

	// max body size
	maxBytes := int64(1048576)

	// parse and transform the comma separated envs that configure CORS
	hs := os.Getenv("CORS_HEADERS")
	ms := os.Getenv("CORS_METHODS")
	ogs := os.Getenv("CORS_ORIGINS")
	normalizeCORSParams := func(e string) []string {
		params := strings.Split(e, ",")
		for i, p := range params {
			params[i] = strings.ReplaceAll(p, " ", "")
		}
		return params
	}
	headers := normalizeCORSParams(hs)
	methods := normalizeCORSParams(ms)
	origins := normalizeCORSParams(ogs)

	// prometheus metric extraction
	ir, ok := promMetrics[MetricKeyInternalRandom].(*prometheus.GaugeVec)
	if !ok {
		return nil, fmt.Errorf("could not find internal random metric")
	}

	// admin/auth handlers
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("GET /ping", handlePing(l, p))
	mux.Handle("POST /token", handleIssueToken(l))

	// workflow schedule routes
	mux.Handle("GET /schedule", stools.AdaptHandler(
		handleGetSchedule(l, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.Handle("POST /schedule", stools.AdaptHandler(
		handleCreateSchedule(l, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.Handle("PUT /schedule", stools.AdaptHandler(
		handleUpdateSchedule(l, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.Handle("DELETE /schedule", stools.AdaptHandler(
		handleCancelSchedule(l, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.Handle("POST /schedule/trigger", stools.AdaptHandler(
		handleTriggerSchedule(l, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))

	// metadata
	mux.HandleFunc("GET /metadata", stools.AdaptHandler(
		handleGetMetricMetadata(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("POST /metadata", stools.AdaptHandler(
		handlePostMetricMetadata(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))

	// internal metrics
	mux.HandleFunc("GET /internal/generate", stools.AdaptHandler(
		handleInternalMetricsGenerate(l),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("GET /internal/metrics", stools.AdaptHandler(
		handleInternalMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		// FIXME: unauthed for now
	))
	mux.HandleFunc("POST /internal/metrics", stools.AdaptHandler(
		handleInternalMetricsPost(l, q, ir),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))

	// youtube video metrics
	mux.HandleFunc("GET /youtube/video", stools.AdaptHandler(
		handleYouTubeMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		// FIXME: unauthed for now
	))
	mux.HandleFunc("POST /youtube/video", stools.AdaptHandler(
		handleYouTubeVideoMetricsPost(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))

	// kaggle notebook metrics
	mux.HandleFunc("GET /kaggle/notebook", stools.AdaptHandler(
		handleKaggleNotebookMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("POST /kaggle/notebook", stools.AdaptHandler(
		handleKaggleNotebookPost(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	// kaggle dataset metrics
	mux.HandleFunc("GET /kaggle/dataset", stools.AdaptHandler(
		handleKaggleDatasetMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("POST /kaggle/dataset", stools.AdaptHandler(
		handleKaggleDatasetPost(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	// reddit post metrics
	mux.HandleFunc("GET /reddit/post", stools.AdaptHandler(
		handleRedditPostMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("POST /reddit/post", stools.AdaptHandler(
		handleRedditPostMetricsPost(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	// reddit comment metrics
	mux.HandleFunc("GET /reddit/comment", stools.AdaptHandler(
		handleRedditCommentMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("POST /reddit/comment", stools.AdaptHandler(
		handleRedditCommentMetricsPost(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))

	// getting timeseries
	mux.HandleFunc("GET /timeseries/raw", stools.AdaptHandler(
		handleGetTimeSeriesByIDs(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("GET /timeseries/binned", stools.AdaptHandler(
		handleGetTimeSeriesByIDsBucketed(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))

	return mux, nil
}
