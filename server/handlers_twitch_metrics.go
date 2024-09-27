package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
	"github.com/prometheus/client_golang/prometheus"
)

func handleTwitchClipMetricsPost(l *slog.Logger, q *dbgen.Queries, pms map[string]prometheus.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.TwitchClipMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}

		// upload metrics
		if p.SetViewCount {
			err = q.InsertTwitchClipViews(
				r.Context(),
				dbgen.InsertTwitchClipViewsParams{
					ID: p.ID, Views: int64(p.ViewCount)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		writeOK(w)
	}
}

func handleTwitchVideoMetricsPost(l *slog.Logger, q *dbgen.Queries, pms map[string]prometheus.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.TwitchClipMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}

		// upload metrics
		if p.SetViewCount {
			err = q.InsertTwitchVideoViews(
				r.Context(),
				dbgen.InsertTwitchVideoViewsParams{
					ID: p.ID, Views: int64(p.ViewCount)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		writeOK(w)
	}
}

func handleTwitchStreamMetricsPost(l *slog.Logger, q *dbgen.Queries, pms map[string]prometheus.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.TwitchClipMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}

		// upload metrics
		if p.SetViewCount {
			err = q.InsertTwitchStreamViews(
				r.Context(),
				dbgen.InsertTwitchStreamViewsParams{
					ID: p.ID, Views: int64(p.ViewCount)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		writeOK(w)
	}
}

func handleTwitchUserPastDecMetricsPost(l *slog.Logger, q *dbgen.Queries, pms map[string]prometheus.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.TwitchUserPastDecMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}

		// upload metrics
		if p.SetAvgViewCount {
			err = q.InsertTwitchUserPastDecAvgViews(
				r.Context(),
				dbgen.InsertTwitchUserPastDecAvgViewsParams{
					ID: p.ID, AvgViews: p.AvgViewCount})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetMedViewCount {
			err = q.InsertTwitchUserPastDecMedViews(
				r.Context(),
				dbgen.InsertTwitchUserPastDecMedViewsParams{
					ID: p.ID, MedViews: p.MedViewCount})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetStdViewCount {
			err = q.InsertTwitchUserPastDecStdViews(
				r.Context(),
				dbgen.InsertTwitchUserPastDecStdViewsParams{
					ID: p.ID, StdViews: p.StdViewCount})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetAvgDuration {
			err = q.InsertTwitchUserPastDecAvgDuration(
				r.Context(),
				dbgen.InsertTwitchUserPastDecAvgDurationParams{
					ID: p.ID, AvgDuration: p.AvgDuration})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetMedDuration {
			err = q.InsertTwitchUserPastDecMedDuration(
				r.Context(),
				dbgen.InsertTwitchUserPastDecMedDurationParams{
					ID: p.ID, MedDuration: p.MedDuration})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetStdDuration {
			err = q.InsertTwitchUserPastDecStdDuration(
				r.Context(),
				dbgen.InsertTwitchUserPastDecStdDurationParams{
					ID: p.ID, StdDuration: p.StdDuration})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		writeOK(w)
	}
}
