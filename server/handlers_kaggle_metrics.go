package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
	"github.com/prometheus/client_golang/prometheus"
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
			// FIXME: I think this needs to check for isPGError(err, noRows)
			writeInternalError(l, w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func handleKaggleNotebookPost(l *slog.Logger, q *dbgen.Queries, votes, downloads *prometheus.GaugeVec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.KaggleMetricPayload
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

		// set vote metrics
		c, err := votes.GetMetricWithLabelValues(p.Slug)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		c.Set(float64(p.Votes))

		// set download metrics
		c, err = downloads.GetMetricWithLabelValues(p.Slug)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		c.Set(float64(p.Downloads))

		if p.SetVotes {
			err = q.InsertKaggleNotebookVotes(r.Context(), dbgen.InsertKaggleNotebookVotesParams{Slug: p.Slug, Votes: int32(p.Votes)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetDownloads {
			err = q.InsertKaggleNotebookDownloads(r.Context(), dbgen.InsertKaggleNotebookDownloadsParams{Slug: p.Slug, Downloads: int32(p.Downloads)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		// wrap up
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
			// FIXME: I think this needs to check for isPGError(err, noRows)
			writeInternalError(l, w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func handleKaggleDatasetPost(l *slog.Logger, q *dbgen.Queries, votes, downloads *prometheus.GaugeVec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p api.KaggleMetricPayload
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

		// set vote metrics
		c, err := votes.GetMetricWithLabelValues(p.Slug)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		c.Set(float64(p.Votes))

		// set download metrics
		c, err = downloads.GetMetricWithLabelValues(p.Slug)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		c.Set(float64(p.Downloads))

		// upload timeseries
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
		// wrap up
		writeOK(w)
	}
}
