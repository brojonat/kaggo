package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
)

func handleRedditPostMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			writeBadRequestError(w, fmt.Errorf("must supply id"))
			return
		}
		res, err := q.GetRedditPostMetricsByID(r.Context(), id)
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

func handleRedditPostMetricsPost(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
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
					ID: p.ID, Title: p.Title, Score: int32(p.Score)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetRatio {
			err = q.InsertRedditPostRatio(
				r.Context(),
				dbgen.InsertRedditPostRatioParams{
					ID: p.ID, Title: p.Title, Ratio: p.Ratio})
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
		id := r.URL.Query().Get("id")
		if id == "" {
			writeBadRequestError(w, fmt.Errorf("must supply id"))
			return
		}
		res, err := q.GetRedditCommentMetricsByID(r.Context(), id)
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

func handleRedditCommentMetricsPost(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
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
