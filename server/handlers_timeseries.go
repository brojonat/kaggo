package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/brojonat/kaggo/server/db/dbgen"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
	"github.com/jackc/pgx/v5/pgtype"
)

func handleGetTimeSeriesByIDsBucketed(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bs := r.URL.Query().Get("bucket_size")
		if bs == "" {
			bs = "15 min"
		}
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id(s)"))
			return
		}
		// FIXME: validate!
		res, err := q.GetYouTubeVideoMetricsByIDsBucketed(
			r.Context(),
			dbgen.GetYouTubeVideoMetricsByIDsBucketedParams{
				Ids:     ids,
				TsStart: pgtype.Timestamptz{Time: time.Time{}, Valid: true},
				TsEnd:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
			},
		)

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

// Return all timeseries metrics under the supplied (kind, id). The return type
// is a list of objects. The objects may not all have the same fields/types,
// since they will have different metrics (views, ratings, etc). It is the
// responsibility of the caller to handle the objects correctly (i.e., they can
// use the "metric" field to infer the respective type).
func handleGetTimeSeriesByIDs(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rk := r.URL.Query().Get("request_kind")
		ids := r.URL.Query()["id"]
		dur := r.URL.Query().Get("dur")
		if rk == "" || len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply request_kind and id(s)"))
			return
		}

		// optionally truncate by timestamp
		var ts_start time.Time
		if dur != "" {
			tdur, err := time.ParseDuration(dur)
			if err != nil {
				writeBadRequestError(w, fmt.Errorf("could not parse duration: %w", err))
				return
			}
			ts_start = time.Now().Add(-tdur)
		}

		var err error
		var rows interface{}

		switch rk {
		case kt.RequestKindInternalRandom:
			rows, err = getInternalRandomTimeSeries(r.Context(), l, q, ids, ts_start, time.Now())
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		case kt.RequestKindYouTubeVideo:
			rows, err = getYouTubeVideoTimeSeries(r.Context(), l, q, ids, ts_start, time.Now())
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		case kt.RequestKindKaggleNotebook:
			rows, err = getKaggleNotebookTimeSeries(r.Context(), l, q, ids, ts_start, time.Now())
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		case kt.RequestKindKaggleDataset:
			rows, err = getKaggleDatasetTimeSeries(r.Context(), l, q, ids, ts_start, time.Now())
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		case kt.RequestKindRedditPost:
			rows, err = getRedditPostTimeSeries(r.Context(), l, q, ids, ts_start, time.Now())
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		case kt.RequestKindRedditComment:
			rows, err = getRedditCommentTimeSeries(r.Context(), l, q, ids, ts_start, time.Now())
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		default:
			writeBadRequestError(w, fmt.Errorf("unexpected RequestKind %s", rk))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rows)

	}
}

// maybe move these into a service file or something

// Main entry point for getting the internal.random timeseries. This will return a
// list of JSON serializable objects. The details are...in a state of flux.
func getInternalRandomTimeSeries(
	ctx context.Context,
	l *slog.Logger,
	q *dbgen.Queries,
	ids []string,
	ts_start time.Time,
	ts_end time.Time,
) (interface{}, error) {
	return q.GetInternalMetricsByIDs(ctx, dbgen.GetInternalMetricsByIDsParams{
		Ids:     ids,
		TsStart: pgtype.Timestamptz{Time: ts_start, Valid: true},
		TsEnd:   pgtype.Timestamptz{Time: ts_end, Valid: true},
	})
}

func getYouTubeVideoTimeSeries(
	ctx context.Context,
	l *slog.Logger,
	q *dbgen.Queries,
	ids []string,
	ts_start time.Time,
	ts_end time.Time,
) (interface{}, error) {
	return q.GetYouTubeVideoMetricsByIDs(ctx, dbgen.GetYouTubeVideoMetricsByIDsParams{
		Ids:     ids,
		TsStart: pgtype.Timestamptz{Time: ts_start, Valid: true},
		TsEnd:   pgtype.Timestamptz{Time: ts_end, Valid: true},
	})
}

func getYouTubeChannelTimeSeries(
	ctx context.Context,
	l *slog.Logger,
	q *dbgen.Queries,
	ids []string,
	ts_start time.Time,
	ts_end time.Time,
) (interface{}, error) {
	return q.GetYouTubeChannelMetricsByIDs(ctx, dbgen.GetYouTubeChannelMetricsByIDsParams{
		Ids:     ids,
		TsStart: pgtype.Timestamptz{Time: ts_start, Valid: true},
		TsEnd:   pgtype.Timestamptz{Time: ts_end, Valid: true},
	})
}

func getKaggleNotebookTimeSeries(
	ctx context.Context,
	l *slog.Logger,
	q *dbgen.Queries,
	ids []string,
	ts_start time.Time,
	ts_end time.Time,
) (interface{}, error) {
	return q.GetKaggleNotebookMetrics(ctx, dbgen.GetKaggleNotebookMetricsParams{
		Ids:     ids,
		TsStart: pgtype.Timestamptz{Time: ts_start, Valid: true},
		TsEnd:   pgtype.Timestamptz{Time: ts_end, Valid: true},
	})
}

func getKaggleDatasetTimeSeries(
	ctx context.Context,
	l *slog.Logger,
	q *dbgen.Queries,
	ids []string,
	ts_start time.Time,
	ts_end time.Time,
) (interface{}, error) {
	return q.GetKaggleDatasetMetrics(ctx, dbgen.GetKaggleDatasetMetricsParams{
		Ids:     ids,
		TsStart: pgtype.Timestamptz{Time: ts_start, Valid: true},
		TsEnd:   pgtype.Timestamptz{Time: ts_end, Valid: true},
	})
}

func getRedditPostTimeSeries(
	ctx context.Context,
	l *slog.Logger,
	q *dbgen.Queries,
	ids []string,
	ts_start time.Time,
	ts_end time.Time,
) (interface{}, error) {
	return q.GetRedditPostMetricsByIDs(ctx, dbgen.GetRedditPostMetricsByIDsParams{
		Ids:     ids,
		TsStart: pgtype.Timestamptz{Time: ts_start, Valid: true},
		TsEnd:   pgtype.Timestamptz{Time: ts_end, Valid: true},
	})
}

func getRedditCommentTimeSeries(
	ctx context.Context,
	l *slog.Logger,
	q *dbgen.Queries,
	ids []string,
	ts_start time.Time,
	ts_end time.Time,
) (interface{}, error) {
	return q.GetRedditCommentMetricsByIDs(ctx, dbgen.GetRedditCommentMetricsByIDsParams{
		Ids:     ids,
		TsStart: pgtype.Timestamptz{Time: ts_start, Valid: true},
		TsEnd:   pgtype.Timestamptz{Time: ts_end, Valid: true},
	})
}

func getRedditSubredditTimeSeries(
	ctx context.Context,
	l *slog.Logger,
	q *dbgen.Queries,
	ids []string,
	ts_start time.Time,
	ts_end time.Time,
) (interface{}, error) {
	return q.GetRedditSubredditMetricsByIDs(ctx, dbgen.GetRedditSubredditMetricsByIDsParams{
		Ids:     ids,
		TsStart: pgtype.Timestamptz{Time: ts_start, Valid: true},
		TsEnd:   pgtype.Timestamptz{Time: ts_end, Valid: true},
	})
}
