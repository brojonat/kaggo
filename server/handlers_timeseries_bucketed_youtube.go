package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/brojonat/kaggo/server/db/dbgen"
	"github.com/jackc/pgx/v5/pgtype"
)

func handleGetYouTubeVideoTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bs := r.URL.Query().Get("bucket_size")
		if bs == "" {
			bs = "60m"
		}
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			idstr := r.URL.Query().Get("ids")
			ids = strings.Split(idstr, ",")
		}
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id(s)"))
			return
		}

		var res interface{}
		var err error

		switch bs {
		case "15m":

			res, err = q.GetYouTubeVideoMetricsByIDsBucket15Min(
				r.Context(),
				dbgen.GetYouTubeVideoMetricsByIDsBucket15MinParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "60m", "1h":
			res, err = q.GetYouTubeVideoMetricsByIDsBucket1Hr(
				r.Context(),
				dbgen.GetYouTubeVideoMetricsByIDsBucket1HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "8h":
			res, err = q.GetYouTubeVideoMetricsByIDsBucket8Hr(
				r.Context(),
				dbgen.GetYouTubeVideoMetricsByIDsBucket8HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "1d":
			res, err = q.GetYouTubeVideoMetricsByIDsBucket1Day(
				r.Context(),
				dbgen.GetYouTubeVideoMetricsByIDsBucket1DayParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		default:
			writeBadRequestError(w, fmt.Errorf("unsupported bucket_size: %s", bs))
			return
		}

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

func handleGetYouTubeChannelTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bs := r.URL.Query().Get("bucket_size")
		if bs == "" {
			bs = "60m"
		}
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			idstr := r.URL.Query().Get("ids")
			ids = strings.Split(idstr, ",")
		}
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id(s)"))
			return
		}

		var res interface{}
		var err error

		switch bs {
		case "15m":

			res, err = q.GetYouTubeChannelMetricsByIDsBucket15Min(
				r.Context(),
				dbgen.GetYouTubeChannelMetricsByIDsBucket15MinParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "60m", "1h":
			res, err = q.GetYouTubeChannelMetricsByIDsBucket1Hr(
				r.Context(),
				dbgen.GetYouTubeChannelMetricsByIDsBucket1HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "8h":
			res, err = q.GetYouTubeChannelMetricsByIDsBucket8Hr(
				r.Context(),
				dbgen.GetYouTubeChannelMetricsByIDsBucket8HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "1d":
			res, err = q.GetYouTubeChannelMetricsByIDsBucket1Day(
				r.Context(),
				dbgen.GetYouTubeChannelMetricsByIDsBucket1DayParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		default:
			writeBadRequestError(w, fmt.Errorf("unsupported bucket_size: %s", bs))
			return
		}

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
