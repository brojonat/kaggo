package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/dbgen"
)

func handleYouTubeMetricsGet(l *slog.Logger, q *dbgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		slug := r.URL.Query().Get("slug")
		if id == "" && slug == "" {
			writeBadRequestError(w, fmt.Errorf("must supply id or slug"))
			return
		}
		if id != "" {
			res, err := q.GetYouTubeVideoMetricsByID(r.Context(), id)
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
			return
		}
		res, err := q.GetYouTubeVideoMetricsBySlug(r.Context(), slug)
		if err != nil {
			writeInternalError(l, w, err)
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
		if p.ID == "" || p.Slug == "" {
			writeBadRequestError(w, fmt.Errorf("must supply id"))
			return
		}

		// FIXME: set prometheus metrics too

		// upload metrics
		if p.SetViews {
			err = q.InsertYouTubeVideoViews(
				r.Context(),
				dbgen.InsertYouTubeVideoViewsParams{
					ID: p.ID, Slug: p.Slug, Views: int32(p.Views)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}
		if p.SetComments {
			err = q.InsertYouTubeVideoComments(
				r.Context(),
				dbgen.InsertYouTubeVideoCommentsParams{
					ID: p.ID, Slug: p.Slug, Comments: int32(p.Comments)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		if p.SetLikes {
			err = q.InsertYouTubeVideoLikes(
				r.Context(),
				dbgen.InsertYouTubeVideoLikesParams{
					ID: p.ID, Slug: p.Slug, Likes: int32(p.Likes)})
			if err != nil {
				writeInternalError(l, w, err)
				return
			}
		}

		// wrap up
		writeOK(w)
	}
}
