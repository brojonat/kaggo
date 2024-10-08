package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"github.com/brojonat/server-tools/stools"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func handleRunMetadataWF(l *slog.Logger, q *dbgen.Queries, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			writeBadRequestError(w, fmt.Errorf("could not read request body: %w", err))
			return
		}
		defer r.Body.Close()
		var body api.GenericScheduleRequestPayload
		err = json.Unmarshal(b, &body)
		if err != nil {
			writeBadRequestError(w, fmt.Errorf("could not parse request body: %w", err))
			return
		}

		// prepare the request to pass to the workflow
		_, serialReq, id, err := makeExternalRequest(q, body.RequestKind, body.ID, true)
		if err != nil {
			if errors.Is(err, errUnsupportedRequestKind) {
				writeBadRequestError(w, fmt.Errorf("%w: %s", err, body.RequestKind))
				return
			}
			writeInternalError(l, w, err)
			return
		}
		// execute a workflow that will fetch the metadata and post it back to the server.
		// this will be a good litmus test for whether or not the client submitted a "good"
		// entity that we can query before the "scheduled" workflow starts running.
		workflowOptions := client.StartWorkflowOptions{
			ID:          id,
			TaskQueue:   os.Getenv("TEMPORAL_TASK_QUEUE"),
			RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: 3},
		}

		_, err = tc.ExecuteWorkflow(
			r.Context(),
			workflowOptions,
			kt.DoMetadataRequestWF,
			kt.DoMetadataRequestWFRequest{RequestKind: body.RequestKind, Serial: serialReq},
		)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		writeOK(w)
	}
}

func handleGetMetricMetadata(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id(s)"))
			return
		}
		res, err := q.GetMetadataByIDs(r.Context(), ids)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		if res == nil {
			writeEmptyResultError(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

// Returns the children for the supplied request_kind and user id. For example,
// returns the child posts for (reddit.user, username).
func handleGetChildrenMetadata(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rk := r.URL.Query().Get("request_kind")
		id := r.URL.Query().Get("id")
		if rk == "" || id == "" {
			writeBadRequestError(w, fmt.Errorf("must supply request_kind and id"))
			return
		}

		var childRK, ownerField string
		switch rk {
		case kt.RequestKindRedditUser:
			childRK = kt.RequestKindRedditPost
			ownerField = "parent_user_name"
		case kt.RequestKindRedditSubreddit:
			childRK = kt.RequestKindRedditPost
			ownerField = "parent_subreddit"
		case kt.RequestKindYouTubeChannel:
			childRK = kt.RequestKindYouTubeVideo
			ownerField = "parent_channel_id"
		default:
			writeBadRequestError(w, fmt.Errorf("unsupported request_kind %s", rk))
			return
		}
		// FIXME: ownerField must be injected somehow
		l.Info("need to inject", "ownerField", ownerField)

		res, err := q.GetChildrenMetadataByID(
			r.Context(),
			dbgen.GetChildrenMetadataByIDParams{
				ID:                id,
				ParentRequestKind: rk,
				ChildRequestKind:  childRK,
			})
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		if res == nil {
			writeEmptyResultError(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func handlePostMetricMetadata(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data api.MetricMetadataPayload
		err := stools.DecodeJSONBody(r, &data)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		err = q.InsertMetadata(
			r.Context(),
			dbgen.InsertMetadataParams{
				ID:          data.ID,
				RequestKind: data.RequestKind,
				Data:        data.Data,
			})
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		writeOK(w)
	}
}

func handleAddListenerSub(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data api.AddListenerSubPayload
		err := stools.DecodeJSONBody(r, &data)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		switch data.RequestKind {
		case kt.RequestKindYouTubeChannel:
			err = q.InsertYouTubeChannelSubscription(r.Context(), data.ID)
		default:
			writeBadRequestError(w, fmt.Errorf("unsupported request_kind %s", data.RequestKind))
			return
		}
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		writeOK(w)
	}
}
