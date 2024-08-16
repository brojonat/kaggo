package temporal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type DoRequestWFRequest struct {
	RequestKind string `json:"request_kind"`
	Serial      []byte `json:"serial"`
}

// DoRequestWF workflow will perform a request to some external service that
// returns metrics (i.e., number of votes for some identifier) and passes the
// response to a metrics handler that will upload them back to our server.
func DoRequestWF(ctx workflow.Context, r DoRequestWFRequest) error {
	logger := workflow.GetLogger(ctx)
	logger.Info(fmt.Sprintf(
		"Starting workflow (wid: %s, rid: %s)",
		workflow.GetInfo(ctx).WorkflowExecution.ID,
		workflow.GetInfo(ctx).WorkflowExecution.RunID))

	var a *ActivityRequester

	// Do the long polling request. Don't retry; these are "cheap" requests and
	// it's better to miss some window of data than risk spamming the external
	// server with retries.
	activityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Second,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 1},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	drp := DoRequestParam(r)
	var drr DoRequestResult
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
		ScheduleToCloseTimeout: 10 * time.Second,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 1},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var urr api.DefaultJSONResponse
	if err := workflow.ExecuteActivity(ctx, a.UploadMetrics, drr).Get(ctx, &urr); err != nil {
		return err
	}
	return nil
}
