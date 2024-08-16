package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
)

func handleKaggleNotebookMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Query().Get("slug")
		if slug == "" {
			writeBadRequestError(w, fmt.Errorf("must supply timeseries slug"))
			return
		}
		res, err := q.GetKaggleNotebookMetrics(r.Context(), slug)
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

func handleKaggleNotebookPost(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.KaggleNotebookMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		if p.Slug == "" {
			writeBadRequestError(w, fmt.Errorf("must supply slug"))
			return
		}

		if p.SetVotes {
			err = q.InsertKaggleNotebookVotes(r.Context(), dbgen.InsertKaggleNotebookVotesParams{Slug: p.Slug, Votes: int32(p.Votes)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		writeOK(w)
	}
}

func handleKaggleDatasetMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Query().Get("slug")
		if slug == "" {
			writeBadRequestError(w, fmt.Errorf("must supply timeseries slug"))
			return
		}
		res, err := q.GetKaggleDatasetMetrics(r.Context(), slug)
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

func handleKaggleDatasetPost(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p api.KaggleDatasetMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		if p.Slug == "" {
			writeBadRequestError(w, fmt.Errorf("must supply slug"))
			return
		}

		if p.SetVotes {
			err = q.InsertKaggleDatasetVotes(
				r.Context(),
				dbgen.InsertKaggleDatasetVotesParams{
					Slug: p.Slug, Votes: int32(p.Votes),
				})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		if p.SetViews {
			err = q.InsertKaggleDatasetViews(
				r.Context(),
				dbgen.InsertKaggleDatasetViewsParams{
					Slug: p.Slug, Views: int32(p.Views),
				})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		if p.SetDownloads {
			err = q.InsertKaggleDatasetDownloads(
				r.Context(),
				dbgen.InsertKaggleDatasetDownloadsParams{
					Slug: p.Slug, Downloads: int32(p.Downloads),
				})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		writeOK(w)
	}
}
