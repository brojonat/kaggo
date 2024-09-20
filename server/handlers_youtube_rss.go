package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"github.com/jackc/pgx/v5"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func handleGetYouTubeWebSubTargets(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, err := q.GetYouTubeChannelSubscriptions(r.Context())
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		body := kt.YouTubeChannelSubActRequest{
			ChannelIDs: ids,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(body)
	}
}

func handleRunYouTubeListener(l *slog.Logger, q *dbgen.Queries, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wopts := client.StartWorkflowOptions{
			ID:          "youtube.listener",
			TaskQueue:   "kaggo",
			RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: 1},
		}
		wfr := kt.RunYouTubeListenerWFRequest{}
		_, err := tc.ExecuteWorkflow(r.Context(), wopts, kt.RunYouTubeListenerWF, wfr)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		writeOK(w)
	}
}

// confirms a websub subscription
func handleYouTubeVideoWebSubSetup(l *slog.Logger, q *dbgen.Queries, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topicRaw := r.URL.Query().Get("hub.topic")
		challenge := r.URL.Query().Get("hub.challenge")
		turl, err := url.Parse(topicRaw)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		cid := turl.Query().Get("channel_id")
		_, err = q.YouTubeChannelSubscriptionExists(r.Context(), cid)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeInternalError(l, w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
	}
}

// handles a YouTube channel WebSub update and makes a schedule for the newly posted video
func handleYouTubeVideoWebSubNotification(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	seen := sync.Map{}
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		xmlstr := string(b)
		re := regexp.MustCompile(`<yt:videoId>(?P<vid>.+)</yt:videoId>`)
		vidIdx := re.SubexpIndex("vid")
		matches := re.FindStringSubmatch(xmlstr)
		if len(matches) < 1 {
			l.Error("could not parse the following xml", "xml", xmlstr)
			writeInternalError(l, w, fmt.Errorf("could not parse video ID"))
			return
		}
		vid := matches[vidIdx]

		// We have a local cache to deal with duplicate notifications. This
		// doesn't have to be perfect, but it'll get us 90% of the way there.
		// THe service will get restarted frequently enough that this shouldn't
		// grow _too_ large.
		if _, ok := seen.Load(vid); ok {
			writeOK(w)
			return
		}
		seen.Store(vid, struct{}{})

		l.Info("got new youtube video to monitor", "id", vid)
		// we want to follow this post for some nominal amount of time
		rk := "youtube.video"
		sched := GetDefaultScheduleSpec(rk, vid)
		sched.EndAt = time.Now().Add(21 * 24 * time.Hour) // 3 weeks
		payload := api.GenericScheduleRequestPayload{
			RequestKind: rk,
			ID:          vid,
			Schedule:    sched,
		}
		b, err = json.Marshal(payload)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		wfr, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("http://localhost:%s", os.Getenv("SERVER_PORT"))+"/schedule",
			bytes.NewReader(b))
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		wfr.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
		res, err := http.DefaultClient.Do(wfr)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		defer res.Body.Close()
		b, err = io.ReadAll(res.Body)
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		// we expect some posts will end up here twice, especially in cases where the
		// creator changes the title or description, so ignore StatusConflict
		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusConflict {
			writeInternalError(l, w, fmt.Errorf("bad response from server: %d: %s", res.StatusCode, string(b)))
			return
		}
		writeOK(w)
	}
}
