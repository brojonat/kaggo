package worker

import (
	"context"
	"log"
	"log/slog"

	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func RunWorker(ctx context.Context, l *slog.Logger, thp string) error {
	// connect to temporal
	c, err := client.Dial(client.Options{
		Logger:   l,
		HostPort: thp,
	})
	if err != nil {
		log.Fatalf("Couldn't initialize Temporal client. Exiting.\nError: %s", err)
	}
	defer c.Close()

	// register workflows
	w := worker.New(c, "kaggo", worker.Options{})
	w.RegisterWorkflow(kt.DoPollingRequestWF)
	w.RegisterWorkflow(kt.DoMetadataRequestWF)

	// register activities
	a := &kt.ActivityRequester{}
	w.RegisterActivity(a)

	// run indefinitely
	return w.Run(worker.InterruptCh())
}
