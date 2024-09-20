package worker

import (
	"context"
	"log"
	"log/slog"
	"os"

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
