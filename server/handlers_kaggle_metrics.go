package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
)

func handleKaggleNotebookMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id(s)"))
			return
		}
		res, err := getKaggleNotebookTimeSeries(r.Context(), l, q, ids, time.Time{}, time.Now())
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
		if p.ID == "" {
			writeBadRequestError(w, fmt.Errorf("must supply id"))
			return
		}

		if p.SetVotes {
			err = q.InsertKaggleNotebookVotes(
				r.Context(),
				dbgen.InsertKaggleNotebookVotesParams{
					ID:    p.ID,
					Votes: int32(p.Votes),
				})
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
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply timeseries slug"))
			return
		}
		res, err := getKaggleDatasetTimeSeries(r.Context(), l, q, ids, time.Time{}, time.Now())
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
		if p.ID == "" {
			writeBadRequestError(w, fmt.Errorf("must supply slug"))
			return
		}

		if p.SetVotes {
			err = q.InsertKaggleDatasetVotes(
				r.Context(),
				dbgen.InsertKaggleDatasetVotesParams{
					ID: p.ID, Votes: int32(p.Votes),
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
					ID: p.ID, Views: int32(p.Views),
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
					ID: p.ID, Downloads: int32(p.Downloads),
				})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		writeOK(w)
	}
}
