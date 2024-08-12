// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: youtube-video-metrics.sql

package dbgen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getYouTubeVideoMetricsByID = `-- name: GetYouTubeVideoMetricsByID :many
SELECT id, slug, ts, views, 'youtube.video.views' AS "metric"
FROM youtube_video_views y
WHERE y.id = $1
UNION ALL
SELECT id, slug, ts, views, 'youtube.video.comments' AS "metric"
FROM youtube_video_views y
WHERE y.id = $1
UNION ALL
SELECT id, slug, ts, likes, 'youtube.video.likes' AS "metric"
FROM youtube_video_likes y
WHERE y.id = $1
`

type GetYouTubeVideoMetricsByIDRow struct {
	ID     string             `json:"id"`
	Slug   string             `json:"slug"`
	Ts     pgtype.Timestamptz `json:"ts"`
	Views  int32              `json:"views"`
	Metric string             `json:"metric"`
}

func (q *Queries) GetYouTubeVideoMetricsByID(ctx context.Context, id string) ([]GetYouTubeVideoMetricsByIDRow, error) {
	rows, err := q.db.Query(ctx, getYouTubeVideoMetricsByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetYouTubeVideoMetricsByIDRow
	for rows.Next() {
		var i GetYouTubeVideoMetricsByIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Slug,
			&i.Ts,
			&i.Views,
			&i.Metric,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getYouTubeVideoMetricsBySlug = `-- name: GetYouTubeVideoMetricsBySlug :many
SELECT id, slug, ts, views, 'youtube.video.views' AS "metric"
FROM youtube_video_views y
WHERE y.slug = $1
UNION ALL
SELECT id, slug, ts, views, 'youtube.video.comments' AS "metric"
FROM youtube_video_views y
WHERE y.slug = $1
UNION ALL
SELECT id, slug, ts, likes, 'youtube.video.likes' AS "metric"
FROM youtube_video_likes y
WHERE y.slug = $1
`

type GetYouTubeVideoMetricsBySlugRow struct {
	ID     string             `json:"id"`
	Slug   string             `json:"slug"`
	Ts     pgtype.Timestamptz `json:"ts"`
	Views  int32              `json:"views"`
	Metric string             `json:"metric"`
}

func (q *Queries) GetYouTubeVideoMetricsBySlug(ctx context.Context, slug string) ([]GetYouTubeVideoMetricsBySlugRow, error) {
	rows, err := q.db.Query(ctx, getYouTubeVideoMetricsBySlug, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetYouTubeVideoMetricsBySlugRow
	for rows.Next() {
		var i GetYouTubeVideoMetricsBySlugRow
		if err := rows.Scan(
			&i.ID,
			&i.Slug,
			&i.Ts,
			&i.Views,
			&i.Metric,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertYouTubeVideoComments = `-- name: InsertYouTubeVideoComments :exec
INSERT INTO youtube_video_comments (id, slug, ts, comments)
VALUES ($1, $2, NOW()::TIMESTAMPTZ, $3)
`

type InsertYouTubeVideoCommentsParams struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Comments int32  `json:"comments"`
}

func (q *Queries) InsertYouTubeVideoComments(ctx context.Context, arg InsertYouTubeVideoCommentsParams) error {
	_, err := q.db.Exec(ctx, insertYouTubeVideoComments, arg.ID, arg.Slug, arg.Comments)
	return err
}

const insertYouTubeVideoLikes = `-- name: InsertYouTubeVideoLikes :exec
INSERT INTO youtube_video_likes (id, slug, ts, likes)
VALUES ($1, $2, NOW()::TIMESTAMPTZ, $3)
`

type InsertYouTubeVideoLikesParams struct {
	ID    string `json:"id"`
	Slug  string `json:"slug"`
	Likes int32  `json:"likes"`
}

func (q *Queries) InsertYouTubeVideoLikes(ctx context.Context, arg InsertYouTubeVideoLikesParams) error {
	_, err := q.db.Exec(ctx, insertYouTubeVideoLikes, arg.ID, arg.Slug, arg.Likes)
	return err
}

const insertYouTubeVideoViews = `-- name: InsertYouTubeVideoViews :exec
INSERT INTO youtube_video_views (id, slug, ts, views)
VALUES ($1, $2, NOW()::TIMESTAMPTZ, $3)
`

type InsertYouTubeVideoViewsParams struct {
	ID    string `json:"id"`
	Slug  string `json:"slug"`
	Views int32  `json:"views"`
}

func (q *Queries) InsertYouTubeVideoViews(ctx context.Context, arg InsertYouTubeVideoViewsParams) error {
	_, err := q.db.Exec(ctx, insertYouTubeVideoViews, arg.ID, arg.Slug, arg.Views)
	return err
}
