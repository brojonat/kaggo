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

const (
	ScheduleKindInternalRandom = "internal.random"
	ScheduleKindYouTubeVideo   = "youtube.video"
	// ScheduleKindKaggleNotebook = "kaggle.notebook"
	// ScheduleKindKaggleDataset = "kaggle.dataset"
	// ScheduleKindRedditPost = "reddit.post"
	// ScheduleKindRedditComment = "reddit.comment"
)

type GenericSchedulePayload struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
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

		// construct request based on ScheduleKind
		var rwf *http.Request
		switch body.Kind {
		case ScheduleKindYouTubeVideo:
			rwf, err = http.NewRequest(http.MethodGet, "https://youtube.googleapis.com/youtube/v3/videos", nil)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
			q := rwf.URL.Query()
			q.Set("part", "snippet,contentDetails,statistics")
			q.Set("key", os.Getenv("YOUTUBE_API_KEY")) // docs indicate this is required in both?
			q.Set("id", body.ID)
			rwf.URL.RawQuery = q.Encode()
			rwf.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("YOUTUBE_API_KEY")))
			rwf.Header.Add("Accept", "application/json")

		default:
			writeBadRequestError(w, fmt.Errorf("unsupported ScheduleKind: %s", body.Kind))
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
		// schedule id is the request url concatenated with the hash of the request
		// whereas the workflow id is simply the url string since it doesn't need to
		// be unique to avoid duplication.
		sid := fmt.Sprintf("%s %x", rwf.URL.String(), h.Sum(nil))
		wid := rwf.URL.String()

		// Create the schedule.
		_, err = tc.ScheduleClient().Create(
			r.Context(),
			client.ScheduleOptions{
				ID: sid,
				Spec: client.ScheduleSpec{
					Calendars: []client.ScheduleCalendarSpec{
						{
							Second:  []client.ScheduleRange{{Start: 0, End: 59, Step: 5}},
							Minute:  []client.ScheduleRange{{Start: 0, End: 59}},
							Hour:    []client.ScheduleRange{{Start: 0, End: 23}},
							Comment: "Every 5 seconds.",
						},
					},
				},
				Action: &client.ScheduleWorkflowAction{
					ID:        wid,
					TaskQueue: "kaggo",
					Workflow:  kt.DoRequestWF,
					Args: []interface{}{kt.DoRequestWFRequest{
						RequestKind: body.Kind, Serial: serialReq}},
				},
			})
		if err != nil {
			if errors.Is(err, temporal.ErrScheduleAlreadyRunning) ||
				stools.IssServiceError(err) {
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

func handleCreateScheduleDeprecated(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint := "https://api.kaggo.brojonat.com/internal/generate"
		rwf, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		rwf.Header.Add("Authorization", r.Header.Get("Authorization"))

		buf := &bytes.Buffer{}
		err = rwf.Write(buf)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		bs := buf.Bytes()

		h := md5.New()
		_, err = h.Write(bs)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}

		sid := rwf.URL.String()
		wid := rwf.URL.String()

		// Create the schedule.
		_, err = tc.ScheduleClient().Create(
			r.Context(),
			client.ScheduleOptions{
				ID: sid,
				Spec: client.ScheduleSpec{
					Calendars: []client.ScheduleCalendarSpec{
						{
							Second: []client.ScheduleRange{{Start: 0, End: 59, Step: 5}},
							Minute: []client.ScheduleRange{{Start: 0, End: 59}},
							Hour:   []client.ScheduleRange{{Start: 0, End: 23}},
						},
					},
				},
				Action: &client.ScheduleWorkflowAction{
					ID:        wid,
					TaskQueue: "kaggo",
					Workflow:  kt.DoRequestWF,
					Args: []interface{}{kt.DoRequestWFRequest{
						ID: r.URL.Query().Get("id"), RequestKind: "internal.random", Serial: bs}},
				},
			})
		if err != nil {
			if errors.Is(err, temporal.ErrScheduleAlreadyRunning) ||
				stools.IssServiceError(err) {
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
