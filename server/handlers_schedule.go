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
	"os"

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
			rwf, err = http.NewRequest(http.MethodGet, "https://api.kaggo.brojonat.com/internal/generate", nil)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
			rwf.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
			rwf.Header.Add("Accept", "application/json")

		case kt.RequestKindYouTubeVideo:
			rwf, err = http.NewRequest(http.MethodGet, "https://youtube.googleapis.com/youtube/v3/videos", nil)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
			q := rwf.URL.Query()
			q.Set("part", "snippet,contentDetails,statistics")
			q.Set("key", os.Getenv("YOUTUBE_API_KEY"))
			q.Set("id", body.ID)
			rwf.URL.RawQuery = q.Encode()
			rwf.Header.Add("Accept", "application/json")

		case kt.RequestKindKaggleNotebook:
			// https://github.com/Kaggle/kaggle-api/blob/48d0433575cac8dd20cf7557c5d749987f5c14a2/kaggle/api/kaggle_api.py#L3052
			rwf, err = http.NewRequest(http.MethodGet, "https://www.kaggle.com/api/v1/kernels/list", nil)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
			// filter to notebooks only and search using the supplied ref
			q := rwf.URL.Query()
			q.Set("search", body.ID)
			rwf.URL.RawQuery = q.Encode()
			// basic auth
			rwf.Header.Add("Accept", "application/json")
			rwf.SetBasicAuth(os.Getenv("KAGGLE_USERNAME"), os.Getenv("KAGGLE_API_KEY"))

		case kt.RequestKindKaggleDataset:
			// https: //github.com/Kaggle/kaggle-api/blob/48d0433575cac8dd20cf7557c5d749987f5c14a2/kaggle/api/kaggle_api.py#L1731
			rwf, err = http.NewRequest(http.MethodGet, "https://www.kaggle.com/api/v1/datasets/list", nil)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
			// search using the supplied ref
			q := rwf.URL.Query()
			q.Set("search", body.ID)
			rwf.URL.RawQuery = q.Encode()
			// basic auth
			rwf.Header.Add("Accept", "application/json")
			rwf.SetBasicAuth(os.Getenv("KAGGLE_USERNAME"), os.Getenv("KAGGLE_API_KEY"))

		case kt.RequestKindRedditPost:
			rwf, err = http.NewRequest(http.MethodGet, "https://reddit.com/api/info.json", nil)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
			q := rwf.URL.Query()
			q.Set("id", fmt.Sprintf("t3_%s", body.ID))
			rwf.URL.RawQuery = q.Encode()
			rwf.Header.Add("Accept", "application/json")
			rwf.Header.Add("User-Agent", "Debian:github.com/brojonat/kaggo/worker:v0.0.1 (by /u/GreaerG)")

		case kt.RequestKindRedditComment:
			rwf, err = http.NewRequest(http.MethodGet, "https://reddit.com/api/info.json", nil)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
			q := rwf.URL.Query()
			q.Set("id", fmt.Sprintf("t1_%s", body.ID))
			rwf.URL.RawQuery = q.Encode()
			rwf.Header.Add("Accept", "application/json")
			rwf.Header.Add("User-Agent", "Debian:github.com/brojonat/kaggo/worker:v0.0.1 (by /u/GreaerG)")

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

		// identifiers are the request kind, identifier, and hash of the request
		id := fmt.Sprintf("%s %s %x", body.RequestKind, body.ID, h.Sum(nil))

		// Create the schedule.
		_, err = tc.ScheduleClient().Create(
			r.Context(),
			client.ScheduleOptions{
				ID:   id,
				Spec: sched,
				Action: &client.ScheduleWorkflowAction{
					ID:        id,
					TaskQueue: "kaggo",
					Workflow:  kt.DoRequestWF,
					Args: []interface{}{kt.DoRequestWFRequest{
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
