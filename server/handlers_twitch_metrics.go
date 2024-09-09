package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
	"github.com/prometheus/client_golang/prometheus"
)

func setTwitchPromMetrics(l *slog.Logger, data api.MetricQueryInternalData, labels prometheus.Labels, pms map[string]prometheus.Collector) {

	// Set Prometheus metrics. The ones for Twitch are here:
	// https://dev.twitch.tv/docs/api/guide/#how-it-works
	mnames := []string{
		PromMetricXRatelimitLimit,
		PromMetricXRatelimitRemaining,
		PromMetricXRatelimitReset,
	}
	for _, mk := range mnames {
		gv, ok := pms[mk].(*prometheus.GaugeVec)
		if !ok {
			l.Error(fmt.Sprintf("failed to locate prom metric %s, skipping", mk))
			continue
		}

		c, err := gv.GetMetricWith(labels)
		if err != nil {
			// GetMetricWith is a get-or-create operation, this should never happen
			l.Error(fmt.Sprintf("failed to get prom metric %s with labels: %s", mk, labels))
			continue
		}

		var val float64

		// NOTE: twitch deviates a little from how these are conventionally supplied,
		// but we're fudging them a bit here to reduce the number of metrics and keep
		// things simple on our end. The main things are that Twitch ratelimit headers
		// don't have the `X-` prefix, and the first one is Limit rather than Used.
		// https://dev.twitch.tv/docs/api/guide/#how-it-works

		switch mk {
		case PromMetricXRatelimitLimit:
			val, err = strconv.ParseFloat(data.RatelimitLimit, 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, data.RatelimitLimit))
				continue
			}
		case PromMetricXRatelimitRemaining:
			val, err = strconv.ParseFloat(data.RatelimitRemaining, 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, data.RatelimitRemaining))
				continue
			}
		case PromMetricXRatelimitReset:
			val, err = strconv.ParseFloat(data.RatelimitReset, 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, data.RatelimitReset))
				continue
			}
		}

		c.Set(val)
	}
}

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

		// prometheus
		labels := prometheus.Labels{"source": "twitch"}
		setTwitchPromMetrics(l, p.InternalData, labels, pms)

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

		// prometheus
		labels := prometheus.Labels{"source": "twitch"}
		setTwitchPromMetrics(l, p.InternalData, labels, pms)

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

		// prometheus
		labels := prometheus.Labels{"source": "twitch"}
		setTwitchPromMetrics(l, p.InternalData, labels, pms)

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

		// prometheus
		labels := prometheus.Labels{"source": "twitch"}
		setTwitchPromMetrics(l, p.InternalData, labels, pms)

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
