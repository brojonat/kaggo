package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
	"github.com/prometheus/client_golang/prometheus"
)

func handleRedditPostMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id"))
			return
		}
		res, err := getRedditPostTimeSeries(r.Context(), l, q, ids, time.Time{}, time.Now())
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

func handleRedditPostMetricsPost(l *slog.Logger, q *dbgen.Queries, pms map[string]prometheus.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.RedditPostMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}

		// upload metrics
		if p.SetScore {
			err = q.InsertRedditPostScore(
				r.Context(),
				dbgen.InsertRedditPostScoreParams{
					ID: p.ID, Score: int32(p.Score)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetRatio {
			err = q.InsertRedditPostRatio(
				r.Context(),
				dbgen.InsertRedditPostRatioParams{
					ID: p.ID, Ratio: p.Ratio})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		writeOK(w)
	}
}

func handleRedditCommentMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id"))
			return
		}
		res, err := getRedditCommentTimeSeries(r.Context(), l, q, ids, time.Time{}, time.Now())
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

func handleRedditCommentMetricsPost(l *slog.Logger, q *dbgen.Queries, pms map[string]prometheus.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.RedditCommentMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}

		// upload metrics
		if p.SetScore {
			err = q.InsertRedditCommentScore(
				r.Context(),
				dbgen.InsertRedditCommentScoreParams{
					ID: p.ID, Score: int32(p.Score)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetControversiality {
			err = q.InsertRedditCommentControversiality(
				r.Context(),
				dbgen.InsertRedditCommentControversialityParams{
					ID: p.ID, Controversiality: p.Controversiality})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		writeOK(w)
	}
}

func handleRedditSubredditMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id"))
			return
		}
		res, err := getRedditSubredditTimeSeries(r.Context(), l, q, ids, time.Time{}, time.Now())
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

func handleRedditSubredditMetricsPost(l *slog.Logger, q *dbgen.Queries, pms map[string]prometheus.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.RedditSubredditMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}

		// upload metrics
		if p.SetSubscribers {
			err = q.InsertRedditSubredditSubscribers(
				r.Context(),
				dbgen.InsertRedditSubredditSubscribersParams{
					ID: p.ID, Subscribers: int32(p.Subscribers)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetActiveUserCount {
			err = q.InsertRedditSubredditActiveUserCount(
				r.Context(),
				dbgen.InsertRedditSubredditActiveUserCountParams{
					ID: p.ID, ActiveUserCount: int32(p.ActiveUserCount)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		writeOK(w)
	}
}

func handleRedditUserMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id"))
			return
		}
		res, err := getRedditUserTimeSeries(r.Context(), l, q, ids, time.Time{}, time.Now())
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

func handleRedditUserMetricsPost(l *slog.Logger, q *dbgen.Queries, pms map[string]prometheus.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.RedditUserMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}

		// upload metrics
		if p.SetAwardeeKarma {
			err = q.InsertRedditUserAwardeeKarma(
				r.Context(),
				dbgen.InsertRedditUserAwardeeKarmaParams{
					ID: p.ID, Karma: int32(p.AwardeeKarma)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetAwarderKarma {
			err = q.InsertRedditUserAwarderKarma(
				r.Context(),
				dbgen.InsertRedditUserAwarderKarmaParams{
					ID: p.ID, Karma: int32(p.AwarderKarma)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetCommentKarma {
			err = q.InsertRedditUserCommentKarma(
				r.Context(),
				dbgen.InsertRedditUserCommentKarmaParams{
					ID: p.ID, Karma: int32(p.CommentKarma)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetLinkKarma {
			err = q.InsertRedditUserLinkKarma(
				r.Context(),
				dbgen.InsertRedditUserLinkKarmaParams{
					ID: p.ID, Karma: int32(p.LinkKarma)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetTotalKarma {
			err = q.InsertRedditUserTotalKarma(
				r.Context(),
				dbgen.InsertRedditUserTotalKarmaParams{
					ID: p.ID, Karma: int32(p.TotalKarma)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		writeOK(w)
	}
}
