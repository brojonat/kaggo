package temporal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Performs a request against an external API and passes the response to a
// handler that will parse metadata from the response and upload the metadata to
// the kaggo server. Typically this will be done once before creating a schedule
// of polling requests to run against this external API.
func DoMetadataRequestWF(ctx workflow.Context, r DoMetadataRequestWFRequest) error {
	var a *ActivityRequester

	// Do a single query to fetch the external data that contains the metadata.
	// Retry this a couple times because we're only doing this once and if
	// we error out then it'll look like a service error.
	activityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 20 * time.Second,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 5},
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

	// Upload the response to our server. Don't retry; if the request doesn't go
	// through, we'll want to fail fast and loud.
	activityOptions = workflow.ActivityOptions{
		ScheduleToCloseTimeout: 20 * time.Second,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 1},
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
		ScheduleToCloseTimeout: 20 * time.Second,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 1},
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

	// Upload the response to our server. Don't retry; if the request doesn't go
	// through, we'll want to fail fast and loud.
	activityOptions = workflow.ActivityOptions{
		ScheduleToCloseTimeout: 20 * time.Second,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 1},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var urr api.DefaultJSONResponse
	umr := UploadMetricsActRequest(drr)
	if err := workflow.ExecuteActivity(ctx, a.UploadMetrics, umr).Get(ctx, &urr); err != nil {
		return err
	}
	return nil
}
