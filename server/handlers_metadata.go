package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brojonat/kaggo/server/db/dbgen"
	"github.com/brojonat/kaggo/server/db/jsonb"
	"github.com/brojonat/server-tools/stools"
)

type MetricMetadataBody struct {
	ID         string             `json:"id"`
	MetricKind string             `json:"metric_kind"`
	Data       jsonb.MetadataJSON `json:"data"`
}

func handleGetMetricMetadata(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id(s)"))
			return
		}
		res, err := q.GetMetadataByIDs(r.Context(), ids)
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

func handlePostMetricMetadata(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data MetricMetadataBody
		err := stools.DecodeJSONBody(r, &data)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}
		err = q.InsertMetadata(
			r.Context(),
			dbgen.InsertMetadataParams{
				ID:         data.ID,
				MetricKind: data.MetricKind,
				Data:       data.Data,
			})
		if err != nil {
			writeInternalError(l, w, err)
			return
		}
		writeOK(w)
	}
}
