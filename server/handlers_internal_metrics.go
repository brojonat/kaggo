package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

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
		id := r.URL.Query().Get("id")
		if id == "" {
			writeBadRequestError(w, fmt.Errorf("must supply timeseries id"))
			return
		}
		res, err := q.GetInternalMetrics(r.Context(), id)
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

func handleInternalMetricsPost(l *slog.Logger, q *dbgen.Queries, value *prometheus.GaugeVec) http.HandlerFunc {
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

		// set vote metrics
		c, err := value.GetMetricWithLabelValues(p.ID)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		c.Set(float64(p.Value))

		err = q.InsertInternalRandom(r.Context(), dbgen.InsertInternalRandomParams{ID: p.ID, Value: int32(p.Value)})
		if err != nil {
			writeInternalError(l, w, err)
			return
		}

		writeOK(w)
	}
}
