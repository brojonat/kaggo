package server

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/brojonat/kaggo/server/api"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"github.com/brojonat/server-tools/stools"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func handleGetSchedule(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ss, err := tc.ScheduleClient().List(r.Context(), client.ScheduleListOptions{})
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		res := []*client.ScheduleListEntry{}
		for {
			if !ss.HasNext() {
				break
			}
			s, err := ss.Next()
			if err != nil {
				break
			}
			res = append(res, s)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

type GenericSchedulePayload struct {
	RequestKind string              `json:"request_kind"`
	ID          string              `json:"id"`
	Schedule    client.ScheduleSpec `json:"schedule_spec"`
}

// create a schedule to query an external api based on the user submitted data
func handleCreateSchedule(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			writeBadRequestError(w, fmt.Errorf("could not read request body: %w", err))
			return
		}
		defer r.Body.Close()
		var body GenericSchedulePayload
		err = json.Unmarshal(b, &body)
		if err != nil {
			writeBadRequestError(w, fmt.Errorf("could not parse request body: %w", err))
			return
		}

		// probably should validate this...but we're the only ones authed for this API and
		// at present we're only using the same fixed schedule, so implement validation
		// later if it's actually needed.
		sched := body.Schedule

		// construct request by switching over RequestKind
		var rwf *http.Request
		switch body.RequestKind {
		case kt.RequestKindInternalRandom:
			rwf, err = makeExternalRequestInternalRandom()
			if err != nil {
				writeInternalError(l, w, err)
				return
			}

		case kt.RequestKindYouTubeVideo:
			rwf, err = makeExternalRequestYouTubeVideo(body.ID)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}

		case kt.RequestKindKaggleNotebook:
			rwf, err = makeExternalRequestKaggleNotebook(body.ID)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}

		case kt.RequestKindKaggleDataset:
			rwf, err = makeExternalRequestKaggleDataset(body.ID)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}

		case kt.RequestKindRedditPost:
			rwf, err = makeExternalRequestRedditPost(body.ID)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}

		case kt.RequestKindRedditComment:
			rwf, err = makeExternalRequestRedditComment(body.ID)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}

		default:
			writeBadRequestError(w, fmt.Errorf("unsupported RequestKind: %s", body.RequestKind))
			return
		}

		// serialize the request
		buf := &bytes.Buffer{}
		err = rwf.Write(buf)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		serialReq := buf.Bytes()
		h := md5.New()
		_, err = h.Write(serialReq)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}

		// identifier is the request kind, identifier, and hash of the request
		id := fmt.Sprintf("%s %s %x", body.RequestKind, body.ID, h.Sum(nil))

		// execute a workflow that will fetch the metadata and post it back to the server.
		// this will be a good litmus test for whether or not the client submitted a "good"
		// entity that we can query before the "scheduled" workflow starts running.
		workflowOptions := client.StartWorkflowOptions{
			ID:          id,
			TaskQueue:   "kaggo",
			RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: 3},
		}

		we, err := tc.ExecuteWorkflow(
			r.Context(),
			workflowOptions,
			kt.DoMetadataRequestWF,
			kt.DoMetadataRequestWFRequest{RequestKind: body.RequestKind, Serial: serialReq},
		)
		// block until this is done; this isn't strictly necessary tbh, once this code
		// is vetted, we can unblock this.
		err = we.Get(r.Context(), &err)
		if err != nil {
			writeInternalError(l, w, fmt.Errorf("error running metadata workflow: %w", err))
			return
		}

		// Create the schedule.
		_, err = tc.ScheduleClient().Create(
			r.Context(),
			client.ScheduleOptions{
				ID:   id,
				Spec: sched,
				Action: &client.ScheduleWorkflowAction{
					ID:        id,
					TaskQueue: "kaggo",
					Workflow:  kt.DoPollingRequestWF,
					Args: []interface{}{kt.DoPollingRequestWFRequest{
						RequestKind: body.RequestKind, Serial: serialReq}},
					RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: 1},
				},
			})
		if err != nil {
			if errors.Is(err, temporal.ErrScheduleAlreadyRunning) ||
				stools.IsTemporalServiceError(err) {
				writeBadRequestError(w, err)
				return
			}
			writeInternalError(l, w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(api.DefaultJSONResponse{Message: "ok"})
	}
}

func handleUpdateSchedule(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		action := r.URL.Query().Get("action")
		if action != "cancel" {
			writeBadRequestError(w, fmt.Errorf("unsupported action: %s", action))
			return
		}
		sid := r.URL.Query().Get("schedule_id")
		note := r.URL.Query().Get("note")
		err := tc.ScheduleClient().GetHandle(r.Context(), sid).Pause(r.Context(), client.SchedulePauseOptions{Note: note})
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		writeOK(w)
	}
}

func handleCancelSchedule(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid := r.URL.Query().Get("schedule_id")
		err := tc.ScheduleClient().GetHandle(r.Context(), sid).Delete(r.Context())
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		writeOK(w)
	}
}

func handleTriggerSchedule(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid := r.URL.Query().Get("schedule_id")
		err := tc.ScheduleClient().GetHandle(r.Context(), sid).Trigger(
			r.Context(),
			client.ScheduleTriggerOptions{Overlap: enums.SCHEDULE_OVERLAP_POLICY_ALLOW_ALL})
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		writeOK(w)
	}
}
