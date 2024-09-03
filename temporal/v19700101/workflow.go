package temporal

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func RunYouTubeListenerWF(ctx workflow.Context, r RunYouTubeListenerWFRequest) error {
	var a *ActivityYouTubeListener

	// get the targets to listen on from the database
	activityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 20 * time.Second,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 1},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var ar YouTubeChannelSubActRequest
	err := workflow.ExecuteActivity(ctx, a.GetYouTubeChannelTargets).Get(ctx, &ar)
	if err != nil {
		return err
	}

	// send all the requests to the websub hub
	activityOptions = workflow.ActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Minute,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 1},
		HeartbeatTimeout:       60 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	return workflow.ExecuteActivity(ctx, a.Subscribe, ar).Get(ctx, nil)
	// FIXME: eventually we'll handle problematic IDs better
}

func RunRedditListenerWF(ctx workflow.Context, r RunRedditListenerWFRequest) error {
	var a *ActivityRedditListener

	// get the targets to listen on from the database
	activityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 20 * time.Second,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 1},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var ar RedditSubActRequest
	err := workflow.ExecuteActivity(ctx, a.GetRedditUserTargets).Get(ctx, &ar)
	if err != nil {
		return err
	}

	// run the long lived monitoring activity
	rp := temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 5.0,
		MaximumInterval:    time.Second * 100,
		MaximumAttempts:    100,
	}
	activityOptions = workflow.ActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Minute,
		RetryPolicy:            &rp,
		HeartbeatTimeout:       60 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	err = workflow.ExecuteActivity(ctx, a.Run, ar).Get(ctx, nil)

	// We expect the activity to eventually timeout; if we don't get a timeout
	// error, then return a generic error, otherwise just restart the workflow
	// as new.
	var te *temporal.TimeoutError
	if !errors.As(err, &te) {
		return fmt.Errorf("unexpected error from activity: %w", err)
	}
	return workflow.NewContinueAsNewError(ctx, RunRedditListenerWF, r)
}

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
