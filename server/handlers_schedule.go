package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func GetDefaultScheduleSpec(rk, id string) client.ScheduleSpec {
	var s client.ScheduleSpec
	switch rk {
	case kt.RequestKindYouTubeChannel, kt.RequestKindYouTubeVideo:
		// do youtube queries every 10 minutes; high res isn't super necessary
		s = client.ScheduleSpec{
			Calendars: []client.ScheduleCalendarSpec{
				{
					Second:  []client.ScheduleRange{{Start: 0}},
					Minute:  []client.ScheduleRange{{Start: 0}},
					Hour:    []client.ScheduleRange{{Start: 0, End: 23, Step: 1}},
					Comment: "every 1 hour",
				},
			},
			Jitter: 60 * 6e9,
		}
	default:
		s = client.ScheduleSpec{
			Calendars: []client.ScheduleCalendarSpec{
				{
					Second:  []client.ScheduleRange{{Start: 0}},
					Minute:  []client.ScheduleRange{{Start: 0, End: 59, Step: 15}},
					Hour:    []client.ScheduleRange{{Start: 0, End: 23}},
					Comment: "every 15 minutes",
				},
			},
			Jitter: 15 * 6e9,
		}
	}
	return s
}

func handleGetSchedule(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rk := r.URL.Query().Get("request_kind")
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
			if rk == "" || strings.HasPrefix(s.ID, rk) {
				res = append(res, s)
			}
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

// create a schedule to query an external api based on the user submitted data
func handleCreateSchedule(l *slog.Logger, q *dbgen.Queries, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		claims, ok := r.Context().Value(ctxKeyJWT).(*authJWTClaims)
		if !ok {
			writeInternalError(l, w, fmt.Errorf("could not extract user email"))
			return
		}

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

		// prepare the request to pass to the metadata workflow
		_, serialReq, id, err := makeExternalRequest(q, body.RequestKind, body.ID, true)
		if err != nil {
			if errors.Is(err, errUnsupportedRequestKind) {
				writeBadRequestError(w, fmt.Errorf("%w: %s", err, body.RequestKind))
				return
			}
			writeInternalError(l, w, err)
			return
		}

		// Execute a workflow that will fetch the metadata and post it back to the server.
		// this will be a good litmus test for whether or not the client submitted a "good"
		// entity that we can query before the "scheduled" workflow starts running.
		workflowOptions := client.StartWorkflowOptions{
			ID:          id,
			TaskQueue:   "kaggo",
			RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: 3},
		}

		// Optionally skip the metadata query. Some clients don't need to run
		// the metadata workflow (e.g., if they're (re)uploading schedules that
		// were deleted). By default, the metadata operation will run and block.
		run_metadata := true
		skip, err := strconv.ParseBool((r.URL.Query().Get("skip-metadata")))
		// ParseBool returns error by default on empty string input, in which
		// case, we should just no-op and stick with running the metadata query.
		if err == nil {
			run_metadata = !skip
		}
		if run_metadata {
			we, err := tc.ExecuteWorkflow(
				r.Context(),
				workflowOptions,
				kt.DoMetadataRequestWF,
				kt.DoMetadataRequestWFRequest{RequestKind: body.RequestKind, Serial: serialReq},
			)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
			err = we.Get(r.Context(), &err)
			if err != nil {
				writeInternalError(l, w, fmt.Errorf("error running metadata workflow: %w", err))
				return
			}
		}

		// probably should validate this...but we're the only ones authed for
		// this API and at present we're only using the same fixed schedule, so
		// implement validation later if it's actually needed.
		sched := body.Schedule

		// prepare the request to pass to the polling workflow
		_, serialReq, id, err = makeExternalRequest(q, body.RequestKind, body.ID, false)
		if err != nil {
			if errors.Is(err, errUnsupportedRequestKind) {
				writeBadRequestError(w, fmt.Errorf("%w: %s", err, body.RequestKind))
				return
			}
			writeInternalError(l, w, err)
			return
		}

		// Create the schedule. Currently we rely on the unique [rk id hash] schedule
		// id to debounce duplicate schedules.
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
			if errors.Is(err, temporal.ErrScheduleAlreadyRunning) {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(api.DefaultJSONResponse{Error: "schedule already running"})
				return
			}
			writeInternalError(l, w, err)
			return
		}

		// add the metric to the user
		p := dbgen.GrantMetricToUserParams{
			Email:       claims.Email,
			RequestKind: body.RequestKind,
			ID:          body.ID,
		}
		if err = q.GrantMetricToUser(r.Context(), p); err != nil {
			l.Error(
				"unable to grant metric to user",
				"email", claims.Email,
				"request_kind", body.RequestKind,
				"id", body.ID,
			)
		}

		// finally, for certain request types, we can opt to monitor the id for
		// new submissions (reddit.users can be monitored for posts and youtube.channels
		// can be monitored for videos).
		if body.Monitor {
			switch body.RequestKind {
			case kt.RequestKindYouTubeChannel:
				err = q.InsertYouTubeChannelSubscription(r.Context(), body.ID)
			case kt.RequestKindRedditUser:
				err = q.InsertYouTubeChannelSubscription(r.Context(), body.ID)
			default:
				err = fmt.Errorf("RequestKind %s doesn't have monitoring support", body.RequestKind)
			}
			if err != nil {
				l.Error(
					"error setting up monitoring",
					"request_kind", body.RequestKind,
					"id", body.ID,
					"error", err.Error(),
				)
			}
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
