package server

import (
	"context"
	"embed"
	"fmt"
	"html/template"
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

//go:embed static
var static embed.FS

// These are prometheus metric keys; handlers may depend on the existence of
// these keys, so any collection, so as a general rule, if you're supplying your
// own prometheus metrics other than the defaults, you should make sure all of
// the following keys are specified.
const (
	PromMetricInternalRandom        = "pm-internal-random"
	PromMetricHandlerRequestCounter = "pm-handler-counter"
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
		PromMetricHandlerRequestCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "request_counter",
				Help: "Count requests to handler",
			},
			[]string{"name", "code"},
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
		HostPort: tcHost,
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

	// max body size, other parsing params
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

	// prometheus setup
	prcounter, ok := pms[PromMetricHandlerRequestCounter].(*prometheus.CounterVec)
	if !ok {
		return nil, fmt.Errorf("error with prom metric %s", PromMetricHandlerRequestCounter)
	}

	// static files and template parsing
	plotTmpl = template.Must(template.ParseFS(static, "static/templates/plots/plot.tmpl"))
	// d3Tmpl = template.Must(template.ParseFS(static, "static/templates/plots/d3.tmpl"))
	// FIXME: get this working so we can iterate
	d3Tmpl = template.Must(template.ParseFiles("server/static/templates/plots/d3.tmpl"))
	// jsFS, err := fs.Sub(static, "static")
	// if err != nil {
	// 	return nil, fmt.Errorf("startup: failed to setup js static file server: %w", err)
	// }
	// fsJS := http.FileServer(http.FS(jsFS))
	fsJS := http.FileServer(http.Dir("server/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fsJS))

	// All route definitions follow; this defines the API surface

	// smoke test/boot handlers
	mux.Handle("GET /ping", stools.AdaptHandler(
		handlePing(l, p),
		withPromCounter(prcounter),
		atLeastOneAuth(
			bearerAuthorizerCtxSetToken(getSecretKey),
			basicAuthorizerCtxSetEmail(getSecretKey),
		),
	))

	// returns a Bearer token; basic auth protected
	mux.Handle("POST /token", stools.AdaptHandler(
		handleIssueToken(l),
		atLeastOneAuth(basicAuthorizerCtxSetEmail(getSecretKey)),
	))
	// serves the static file that will asynchronously fetch the plot data;
	// basic auth protected
	mux.Handle("GET /plots", stools.AdaptHandler(
		handleGetPlots(l, q),
		atLeastOneAuth(basicAuthorizerCtxSetEmail(getSecretKey)),
	))
	// serves the data to the plot static file; Bearer token protected
	mux.Handle("GET /plot-data", stools.AdaptHandler(
		handleGetPlotData(l, q, tc),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
	))
	// prometheus metric handler
	mux.Handle("/metrics", stools.AdaptHandler(
		handlePromMetrics(promhttp.Handler()),
		atLeastOneAuth(
			bearerAuthorizerCtxSetToken(getSecretKey),
			basicAuthorizerCtxSetEmail(getSecretKey),
		),
	))

	// users
	mux.HandleFunc("GET /users", stools.AdaptHandler(
		handleGetUsers(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /users", stools.AdaptHandler(
		handleAddUser(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("DELETE /users", stools.AdaptHandler(
		handleDeleteUser(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("GET /users/metrics", stools.AdaptHandler(
		handleGetUserMetrics(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /users/metrics", stools.AdaptHandler(
		handleUserMetricOperation(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// metadata by metrics
	mux.HandleFunc("GET /metadata", stools.AdaptHandler(
		handleGetMetricMetadata(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("GET /metadata/children", stools.AdaptHandler(
		handleGetChildrenMetadata(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /metadata", stools.AdaptHandler(
		handlePostMetricMetadata(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /metadata/run-workflow", stools.AdaptHandler(
		handleRunMetadataWF(l, q, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// listener subscriptions
	mux.HandleFunc("POST /add-listener-sub", stools.AdaptHandler(
		handleAddListenerSub(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// workflow schedule routes
	mux.Handle("GET /schedule", stools.AdaptHandler(
		handleGetSchedule(l, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.Handle("POST /schedule", stools.AdaptHandler(
		handleCreateSchedule(l, q, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.Handle("PUT /schedule", stools.AdaptHandler(
		handleUpdateSchedule(l, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.Handle("DELETE /schedule", stools.AdaptHandler(
		handleCancelSchedule(l, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.Handle("POST /schedule/trigger", stools.AdaptHandler(
		handleTriggerSchedule(l, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// internal metrics
	mux.HandleFunc("GET /internal/generate", stools.AdaptHandler(
		handleInternalMetricsGenerate(l),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("GET /internal/metrics", stools.AdaptHandler(
		handleInternalMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /internal/metrics", stools.AdaptHandler(
		handleInternalMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// kaggle notebook metrics
	mux.HandleFunc("GET /kaggle/notebook", stools.AdaptHandler(
		handleKaggleNotebookMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /kaggle/notebook", stools.AdaptHandler(
		handleKaggleNotebookPost(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// kaggle dataset metrics
	mux.HandleFunc("GET /kaggle/dataset", stools.AdaptHandler(
		handleKaggleDatasetMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /kaggle/dataset", stools.AdaptHandler(
		handleKaggleDatasetPost(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// youtube video metrics
	mux.HandleFunc("GET /youtube/video", stools.AdaptHandler(
		handleYouTubeVideoMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /youtube/video", stools.AdaptHandler(
		handleYouTubeVideoMetricsPost(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// youtube channel metrics
	mux.HandleFunc("GET /youtube/channel", stools.AdaptHandler(
		handleYouTubeChannelMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /youtube/channel", stools.AdaptHandler(
		handleYouTubeChannelMetricsPost(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// reddit post metrics
	mux.HandleFunc("GET /reddit/post", stools.AdaptHandler(
		handleRedditPostMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /reddit/post", stools.AdaptHandler(
		handleRedditPostMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// reddit comment metrics
	mux.HandleFunc("GET /reddit/comment", stools.AdaptHandler(
		handleRedditCommentMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /reddit/comment", stools.AdaptHandler(
		handleRedditCommentMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// reddit subreddit metrics
	mux.HandleFunc("GET /reddit/subreddit", stools.AdaptHandler(
		handleRedditSubredditMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /reddit/subreddit", stools.AdaptHandler(
		handleRedditSubredditMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// reddit user metrics
	mux.HandleFunc("GET /reddit/user", stools.AdaptHandler(
		handleRedditUserMetricsGet(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /reddit/user", stools.AdaptHandler(
		handleRedditUserMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// twitch clip metrics
	mux.HandleFunc("POST /twitch/clip", stools.AdaptHandler(
		handleTwitchClipMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// twitch video metrics
	mux.HandleFunc("POST /twitch/video", stools.AdaptHandler(
		handleTwitchVideoMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// twitch stream metrics
	mux.HandleFunc("POST /twitch/stream", stools.AdaptHandler(
		handleTwitchStreamMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// twitch user-past-dec metrics
	mux.HandleFunc("POST /twitch/user-past-dec", stools.AdaptHandler(
		handleTwitchUserPastDecMetricsPost(l, q, pms),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// getting timeseries
	mux.HandleFunc("GET /timeseries/raw", stools.AdaptHandler(
		handleGetTimeSeriesByIDs(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("GET /timeseries/bucketed", stools.AdaptHandler(
		handleGetTimeSeriesByIDsBucketed(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	// youtube notifications
	mux.HandleFunc("GET /notification/youtube/targets", stools.AdaptHandler(
		handleGetYouTubeWebSubTargets(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("GET /notification/youtube/websub", stools.AdaptHandler(
		handleYouTubeVideoWebSubSetup(l, q, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /notification/youtube/websub", stools.AdaptHandler(
		handleYouTubeVideoWebSubNotification(l, q),
		apiMode(l, maxBytes, headers, methods, origins),
		withPromCounter(prcounter),
	))
	mux.HandleFunc("POST /run-youtube-listener-wf", stools.AdaptHandler(
		handleRunYouTubeListener(l, q, tc),
		apiMode(l, maxBytes, headers, methods, origins),
		atLeastOneAuth(bearerAuthorizerCtxSetToken(getSecretKey)),
		withPromCounter(prcounter),
	))

	return mux, nil
}
