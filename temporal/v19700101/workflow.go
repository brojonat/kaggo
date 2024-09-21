package temporal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func RunYouTubeListenerWF(ctx workflow.Context, r RunYouTubeListenerWFRequest) error {
	var a *ActivityYouTubeListener

	// Get the targets to listen on from the database. This could fail if we
	// happen to be redeploying; this should retry a bunch
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 20 * time.Second,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 20},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var ar YouTubeChannelSubActRequest
	err := workflow.ExecuteActivity(ctx, a.GetYouTubeChannelTargets).Get(ctx, &ar)
	if err != nil {
		return err
	}

	// Send all the requests to the websub hub. This should be fairly quick because
	// it's just sending a bunch of requests to the pubsubhubbub endpoint.
	activityOptions = workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 5},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	return workflow.ExecuteActivity(ctx, a.Subscribe, ar).Get(ctx, nil)
	// FIXME: eventually we'll handle the case where some IDs were successfully
	// subscribed to and some were not, but that seems like a rare edge case.
}

// Performs a request against an external API and passes the response to a
// handler that will parse metadata from the response and upload the metadata to
// the kaggo server. Typically this will be done once before creating a schedule
// of polling requests to run against this external API.
func DoMetadataRequestWF(ctx workflow.Context, r DoMetadataRequestWFRequest) error {

	var a *ActivityRequester

	// Do a single query to fetch the external data that contains the metadata.
	// Retry this a couple times because we're only doing this once.
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 10, BackoffCoefficient: 5},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	drp := DoRequestActRequest(r)
	var drr DoRequestActResult
	if err := workflow.ExecuteActivity(ctx, a.DoRequest, drp).Get(ctx, &drr); err != nil {
		return err
	}
	if drr.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response: %d (%s): %s",
			drr.StatusCode, http.StatusText(drr.StatusCode), drr.Body)
	}

	// Upload the response to our server. This should also have a bunch of
	// retries associated with it, since we're only doing this once and our
	// server could be down for any number of reasons.
	activityOptions = workflow.ActivityOptions{
		StartToCloseTimeout: 1. * time.Minute,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 10, BackoffCoefficient: 5},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var urr api.DefaultJSONResponse
	umr := UploadMetadataActRequest(drr)
	if err := workflow.ExecuteActivity(ctx, a.UploadMetadata, umr).Get(ctx, &urr); err != nil {
		return err
	}
	return nil
}

// DoPollingRequestWF workflow performs a request against some external API and
// passes the response to a handler that parses metrics from the response and
// uploads the metrics to the kaggo server.
func DoPollingRequestWF(ctx workflow.Context, r DoPollingRequestWFRequest) error {

	var a *ActivityRequester

	// Do the long polling request. Don't retry; these are "cheap" requests and
	// it's better to miss some window of data than risk spamming the external
	// server with retries.
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 1},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	drp := DoRequestActRequest(r)
	var drr DoRequestActResult
	if err := workflow.ExecuteActivity(ctx, a.DoRequest, drp).Get(ctx, &drr); err != nil {
		return err
	}
	if drr.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response: %d (%s): %s",
			drr.StatusCode, http.StatusText(drr.StatusCode), drr.Body)
	}

	// Upload the response to our server. We can retry a couple times over a
	// minute or two in case our server is offline/rebooting, but there's no
	// need to spend a lot of time on this. There will be another polling loop
	// anyway, we can drop some data points here and there. Note that the
	// StartToCloseTimeout needs to accommodate the "monitor" request kinds,
	// which (concurrently) creates a bunch of schedules; this takes a bit
	// because schedule creation is blocked by the initial metadata workflow.
	activityOptions = workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 5, BackoffCoefficient: 5},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var urr api.DefaultJSONResponse
	umr := UploadMetricsActRequest(drr)
	if err := workflow.ExecuteActivity(ctx, a.UploadMetrics, umr).Get(ctx, &urr); err != nil {
		return err
	}
	return nil
}
