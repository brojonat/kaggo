package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
	"github.com/prometheus/client_golang/prometheus"
)

func handleInternalMetricsGenerate(l *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			ID    string `json:"id"`
			Value int    `json:"value"`
		}{
			ID:    "internal-random-metric-identifier",
			Value: rand.Intn(1000),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func handleInternalMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply timeseries id"))
			return
		}
		res, err := getInternalRandomTimeSeries(r.Context(), l, q, ids, time.Time{}, time.Now())
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

func handleInternalMetricsPost(l *slog.Logger, q *dbgen.Queries, pms map[string]prometheus.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.InternalMetricPayload
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

		// Set Prometheus metrics. This is a dummy example that simply records
		// the value associated with the supplied ID.
		gv, ok := pms[PromMetricInternalRandom].(*prometheus.GaugeVec)
		if ok {
			labels := prometheus.Labels{"id": p.ID}
			c, err := gv.GetMetricWith(labels)
			if err != nil {
				// GetMetricWith is a get-or-create operation, this should never happen
				l.Error("failed to get prometheus metric with labels \"id\": %s, skipping", p.ID)
			} else {
				c.Set(float64(p.Value))
			}
		} else {
			l.Error("failed to locate PromMetricInternalRandom, skipping")
		}

		err = q.InsertInternalRandom(
			r.Context(),
			dbgen.InsertInternalRandomParams{
				ID:    p.ID,
				Value: int32(p.Value),
			})
		if err != nil {
			writeInternalError(l, w, err)
			return
		}

		writeOK(w)
	}
}
