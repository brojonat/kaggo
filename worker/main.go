package worker

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	kt "github.com/brojonat/kaggo/temporal/v19700101"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/client"
	sdktally "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/worker"
)

func RunWorker(ctx context.Context, l *slog.Logger, thp string) error {
	// connect to temporal
	c, err := client.Dial(client.Options{
		Logger:   l,
		HostPort: thp,
		MetricsHandler: sdktally.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
			ListenAddress: "0.0.0.0:9090",
			TimerType:     "histogram",
		})),
	})
	if err != nil {
		log.Fatalf("Couldn't initialize Temporal client. Exiting.\nError: %s", err)
	}
	defer c.Close()

	// register workflows
	w := worker.New(c, os.Getenv("TEMPORAL_TASK_QUEUE"), worker.Options{})
	w.RegisterWorkflow(kt.DoPollingRequestWF)
	w.RegisterWorkflow(kt.DoMetadataRequestWF)
	w.RegisterWorkflow(kt.RunYouTubeListenerWF)

	// register activities
	// NOTE: you MUST NOT have any identical methods on these activity structs,
	// or you will encounter a runtime error that prevents all of your workers
	// from starting :O
	a := &kt.ActivityRequester{}
	ysub := &kt.ActivityYouTubeListener{}
	w.RegisterActivity(a)
	w.RegisterActivity(ysub)
	return w.Run(worker.InterruptCh())

}

func newPrometheusScope(c prometheus.Configuration) tally.Scope {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				log.Println("error in prometheus reporter", err)
			},
		},
	)
	if err != nil {
		log.Fatalln("error creating prometheus reporter", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
		Prefix:          "temporal_samples",
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)

	log.Println("prometheus metrics scope created")
	return scope
}
