package server

import (
	"log/slog"
	"net/http"

	"github.com/brojonat/kaggo/server/db/dbgen"
)

func handleMetricsPost(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
