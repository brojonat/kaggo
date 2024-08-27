package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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

func setRedditPromMetrics(l *slog.Logger, data api.MetricQueryInternalData, labels prometheus.Labels, pms map[string]prometheus.Collector) {

	// Set Prometheus metrics. The ones we're interested in for Reddit are
	// the X-Requestlimit-* header values. Range over metrics and set them.
	mnames := []string{
		PromMetricXRatelimitUsed,
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

		switch mk {
		case PromMetricXRatelimitUsed:
			val, err = strconv.ParseFloat(data.XRatelimitUsed, 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, data.XRatelimitUsed))
				continue
			}
		case PromMetricXRatelimitRemaining:
			val, err = strconv.ParseFloat(data.XRatelimitRemaining, 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, data.XRatelimitRemaining))
				continue
			}
		case PromMetricXRatelimitReset:
			val, err = strconv.ParseFloat(data.XRatelimitReset, 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, data.XRatelimitReset))
				continue
			}
		}

		c.Set(val)
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

		// prometheus
		labels := prometheus.Labels{}
		setRedditPromMetrics(l, p.InternalData, labels, pms)

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

		// prometheus
		labels := prometheus.Labels{}
		setRedditPromMetrics(l, p.InternalData, labels, pms)

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

		// prometheus
		labels := prometheus.Labels{}
		setRedditPromMetrics(l, p.InternalData, labels, pms)

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
