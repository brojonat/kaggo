package temporal

import (
	"fmt"
	"net/http"
	"time"

	"go.temporal.io/sdk/workflow"
)

type DoRequestWFRequest struct {
	ID          string `json:"id"`
	RequestKind string `json:"request_kind"`
	Serial      []byte `json:"serial"`
}
type DoRequestWFResponse struct {
	StatusCode int    `json:"status_code"`
	Body       []byte `json:"body"`
}

// DoRequest workflow will perform a request to some external service that
// returns metrics (i.e., number of votes for some identifier) and passes the
// response to a metrics handler that will upload them back to our server.
func DoRequestWF(ctx workflow.Context, r DoRequestWFRequest) error {
	logger := workflow.GetLogger(ctx)
	logger.Info(fmt.Sprintf("Starting workflow ID %s", r.ID))

	// do the long polling request
	activityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	drp := DoRequestParam{
		Serial: r.Serial,
	}
	var drr DoRequestResult
	if err := workflow.ExecuteActivity(ctx, DoRequest, drp).Get(ctx, &drr); err != nil {
		return err
	}
	if drr.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response: %d (%s): %s",
			drr.StatusCode, http.StatusText(drr.StatusCode), drr.Body)
	}
	// take the output and upload the result to our timeseries
	activityOptions = workflow.ActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var urr UploadResponseResult
	if err := workflow.ExecuteActivity(ctx, UploadMetrics, r.RequestKind, drr.Body).Get(ctx, &urr); err != nil {
		return err
	}
	return nil
}
