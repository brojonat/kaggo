package server

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/brojonat/kaggo/server/db/dbgen"
	"github.com/brojonat/server-tools/stools"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.temporal.io/sdk/client"
)

// These are prometheus metric keys; handlers may depend on the existence of
// these keys, so any collection, so as a general rule, if you're supplying your
// own prometheus metrics other than the defaults, you should make sure all of
// the following keys are specified.
const (
	PromMetricInternalRandom      = "pm-internal-random"
	PromMetricXRatelimitUsed      = "pm-x-ratelimit-used"
	PromMetricXRatelimitRemaining = "pm-x-ratelimit-remaining"
	PromMetricXRatelimitReset     = "pm-x-ratelimit-reset"
)

// This is a convenience method for getting the necessary metrics. Some handlers
// (e.g., the internal dummy handler, as well as the reddit handlers) assume
// that particular metrics are passed in as particular types. The handlers will
// log an error and no-op in the event of misconfigured Prometheus metrics, but
// nevertheless you should still probably use the default values provided here
// unless you really know what you're doing.
func GetDefaultPromMetrics() map[string]prometheus.Collector {
	return map[string]prometheus.Collector{
		PromMetricInternalRandom: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "internal_random",
				Help: "A pseudo random metric",
			},
			[]string{"id"},
		),
		PromMetricXRatelimitUsed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "x_ratelimit_used",
				Help: "The X-Ratelimit-Used header from an external server.",
			},
			[]string{"id", "request_kind"},
		),
		PromMetricXRatelimitRemaining: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "x_ratelimit_remaining",
				Help: "The X-Ratelimit-Remaining header from an external server.",
			},
			[]string{"id", "request_kind"},
		),
		PromMetricXRatelimitReset: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "x_ratelimit_reset",
				Help: "The X-Ratelimit-Reset header from an external server.",
			},
			[]string{"id", "request_kind"},
		),
	}
}

// Run the HTTP server.
//
// A note on Prometheus metrics: some handlers expect certain Prometheus metrics to be
// passed in. In the event that you forget, the server will log an error and
// simply no-op on setting the metric, but as a general rule, if you're
// supplying your own prometheus metrics other than those returned by
// GetDefaultPromMetrics, you should make sure all of the PromMetric* keys
// listed above have a corresponding metric passed in.
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

	prometheus.MustRegister(slices.Collect(maps.Values(promMetrics))...)
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
	pms map[string]prometheus.Collector,
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

	// admin/auth handlers
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("GET /ping", handlePing(l, p))
	mux.Handle("POST /token", handleIssueToken())

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
	mux.HandleFunc("POST /metadata/run-workflow", stools.AdaptHandler(
		handleRunMetadataWF(l, tc),
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
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("POST /internal/metrics", stools.AdaptHandler(
		handleInternalMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))

	// youtube video metrics
	mux.HandleFunc("GET /youtube/video", stools.AdaptHandler(
		handleYouTubeVideoMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("POST /youtube/video", stools.AdaptHandler(
		handleYouTubeVideoMetricsPost(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))

	// youtube channel metrics
	mux.HandleFunc("GET /youtube/channel", stools.AdaptHandler(
		handleYouTubeChannelMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("POST /youtube/channel", stools.AdaptHandler(
		handleYouTubeChannelMetricsPost(l, q),
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
		handleRedditPostMetricsPost(l, q, pms),
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
		handleRedditCommentMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	// reddit subreddit metrics
	mux.HandleFunc("GET /reddit/subreddit", stools.AdaptHandler(
		handleRedditSubredditMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("POST /reddit/subreddit", stools.AdaptHandler(
		handleRedditSubredditMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	// reddit subreddit metrics
	mux.HandleFunc("GET /reddit/subreddit", stools.AdaptHandler(
		handleRedditSubredditMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizer(getSecretKey)),
	))
	mux.HandleFunc("POST /reddit/subreddit", stools.AdaptHandler(
		handleRedditSubredditMetricsPost(l, q, pms),
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
