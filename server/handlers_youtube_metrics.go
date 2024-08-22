package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
)

func handleYouTubeMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["id"]
		if len(ids) == 0 {
			writeBadRequestError(w, fmt.Errorf("must supply id(s)"))
			return
		}
		res, err := getYouTubeVideoTimeSeries(r.Context(), l, q, ids, time.Time{}, time.Now())
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

func handleYouTubeVideoMetricsPost(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse
		var p api.YouTubeVideoMetricPayload
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeBadRequestError(w, err)
			return
		}

		// upload metrics
		if p.SetViews {
			err = q.InsertYouTubeVideoViews(
				r.Context(),
				dbgen.InsertYouTubeVideoViewsParams{
					ID:    p.ID,
					Views: int32(p.Views),
				})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetComments {
			err = q.InsertYouTubeVideoComments(
				r.Context(),
				dbgen.InsertYouTubeVideoCommentsParams{
					ID:       p.ID,
					Comments: int32(p.Comments),
				})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		if p.SetLikes {
			err = q.InsertYouTubeVideoLikes(
				r.Context(),
				dbgen.InsertYouTubeVideoLikesParams{
					ID:    p.ID,
					Likes: int32(p.Likes),
				})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		writeOK(w)
	}
}
