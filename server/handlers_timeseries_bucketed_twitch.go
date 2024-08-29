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

func handleGetTwitchClipTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
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

			res, err = q.GetTwitchClipMetricsByIDsBucket15Min(
				r.Context(),
				dbgen.GetTwitchClipMetricsByIDsBucket15MinParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "60m", "1h":
			res, err = q.GetTwitchClipMetricsByIDsBucket1Hr(
				r.Context(),
				dbgen.GetTwitchClipMetricsByIDsBucket1HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "8h":
			res, err = q.GetTwitchClipMetricsByIDsBucket8Hr(
				r.Context(),
				dbgen.GetTwitchClipMetricsByIDsBucket8HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "1d":
			res, err = q.GetTwitchClipMetricsByIDsBucket1Day(
				r.Context(),
				dbgen.GetTwitchClipMetricsByIDsBucket1DayParams{
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

func handleGetTwitchVideoTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
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

			res, err = q.GetTwitchVideoMetricsByIDsBucket15Min(
				r.Context(),
				dbgen.GetTwitchVideoMetricsByIDsBucket15MinParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "60m", "1h":
			res, err = q.GetTwitchVideoMetricsByIDsBucket1Hr(
				r.Context(),
				dbgen.GetTwitchVideoMetricsByIDsBucket1HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "8h":
			res, err = q.GetTwitchVideoMetricsByIDsBucket8Hr(
				r.Context(),
				dbgen.GetTwitchVideoMetricsByIDsBucket8HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "1d":
			res, err = q.GetTwitchVideoMetricsByIDsBucket1Day(
				r.Context(),
				dbgen.GetTwitchVideoMetricsByIDsBucket1DayParams{
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

func handleGetTwitchStreamTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
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

			res, err = q.GetTwitchStreamMetricsByIDsBucket15Min(
				r.Context(),
				dbgen.GetTwitchStreamMetricsByIDsBucket15MinParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "60m", "1h":
			res, err = q.GetTwitchStreamMetricsByIDsBucket1Hr(
				r.Context(),
				dbgen.GetTwitchStreamMetricsByIDsBucket1HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "8h":
			res, err = q.GetTwitchStreamMetricsByIDsBucket8Hr(
				r.Context(),
				dbgen.GetTwitchStreamMetricsByIDsBucket8HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "1d":
			res, err = q.GetTwitchStreamMetricsByIDsBucket1Day(
				r.Context(),
				dbgen.GetTwitchStreamMetricsByIDsBucket1DayParams{
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

func handleGetTwitchUserPastDecTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
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

			res, err = q.GetTwitchUserPastDecMetricsByIDsBucket15Min(
				r.Context(),
				dbgen.GetTwitchUserPastDecMetricsByIDsBucket15MinParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "60m", "1h":
			res, err = q.GetTwitchUserPastDecMetricsByIDsBucket1Hr(
				r.Context(),
				dbgen.GetTwitchUserPastDecMetricsByIDsBucket1HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "8h":
			res, err = q.GetTwitchUserPastDecMetricsByIDsBucket8Hr(
				r.Context(),
				dbgen.GetTwitchUserPastDecMetricsByIDsBucket8HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "1d":
			res, err = q.GetTwitchUserPastDecMetricsByIDsBucket1Day(
				r.Context(),
				dbgen.GetTwitchUserPastDecMetricsByIDsBucket1DayParams{
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
