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

func handleGetRedditPostTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
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

			res, err = q.GetRedditPostMetricsByIDsBucket15Min(
				r.Context(),
				dbgen.GetRedditPostMetricsByIDsBucket15MinParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "60m", "1h":
			res, err = q.GetRedditPostMetricsByIDsBucket1Hr(
				r.Context(),
				dbgen.GetRedditPostMetricsByIDsBucket1HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "8h":
			res, err = q.GetRedditPostMetricsByIDsBucket8Hr(
				r.Context(),
				dbgen.GetRedditPostMetricsByIDsBucket8HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "1d":
			res, err = q.GetRedditPostMetricsByIDsBucket1Day(
				r.Context(),
				dbgen.GetRedditPostMetricsByIDsBucket1DayParams{
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

func handleGetRedditCommentTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
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

			res, err = q.GetRedditCommentMetricsByIDsBucket15Min(
				r.Context(),
				dbgen.GetRedditCommentMetricsByIDsBucket15MinParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "60m", "1h":
			res, err = q.GetRedditCommentMetricsByIDsBucket1Hr(
				r.Context(),
				dbgen.GetRedditCommentMetricsByIDsBucket1HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "8h":
			res, err = q.GetRedditCommentMetricsByIDsBucket8Hr(
				r.Context(),
				dbgen.GetRedditCommentMetricsByIDsBucket8HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "1d":
			res, err = q.GetRedditCommentMetricsByIDsBucket1Day(
				r.Context(),
				dbgen.GetRedditCommentMetricsByIDsBucket1DayParams{
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

func handleGetRedditSubredditTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
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

			res, err = q.GetRedditSubredditMetricsByIDsBucket15Min(
				r.Context(),
				dbgen.GetRedditSubredditMetricsByIDsBucket15MinParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "60m", "1h":
			res, err = q.GetRedditSubredditMetricsByIDsBucket1Hr(
				r.Context(),
				dbgen.GetRedditSubredditMetricsByIDsBucket1HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "8h":
			res, err = q.GetRedditSubredditMetricsByIDsBucket8Hr(
				r.Context(),
				dbgen.GetRedditSubredditMetricsByIDsBucket8HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "1d":
			res, err = q.GetRedditSubredditMetricsByIDsBucket1Day(
				r.Context(),
				dbgen.GetRedditSubredditMetricsByIDsBucket1DayParams{
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

func handleGetRedditUserTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse bucket_size, default to 1 hour
		bs := r.URL.Query().Get("bucket_size")
		if bs == "" {
			bs = "60m"
		}
		// support both id=1&id=2 as well as ids=1,2
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

			res, err = q.GetRedditUserMetricsByIDsBucket15Min(
				r.Context(),
				dbgen.GetRedditUserMetricsByIDsBucket15MinParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "60m", "1h":
			res, err = q.GetRedditUserMetricsByIDsBucket1Hr(
				r.Context(),
				dbgen.GetRedditUserMetricsByIDsBucket1HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "8h":
			res, err = q.GetRedditUserMetricsByIDsBucket8Hr(
				r.Context(),
				dbgen.GetRedditUserMetricsByIDsBucket8HrParams{
					Ids:     ids,
					TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
					TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			)

		case "1d":
			res, err = q.GetRedditUserMetricsByIDsBucket1Day(
				r.Context(),
				dbgen.GetRedditUserMetricsByIDsBucket1DayParams{
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
