package main

import (
	"log"
	"os"

	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// connect to temporal
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_HOST"),
	})
	if err != nil {
		log.Fatalf("Couldn't initialize Temporal client. Exiting.\nError: %s", err)
	}
	defer c.Close()

	// register workflows
	w := worker.New(c, "kaggo", worker.Options{})
	w.RegisterWorkflow(kt.DoRequestWF)

	// register activities
	w.RegisterActivity(kt.DoRequest)
	w.RegisterActivity(kt.UploadMetrics)

	// run indefinitely
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

}
