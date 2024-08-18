// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: youtube-video-metrics.sql

package dbgen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getYouTubeVideoMetricsByIDs = `-- name: GetYouTubeVideoMetricsByIDs :many
SELECT
    y.id AS "id",
    y.title AS "title",
    y.ts AS "ts",
    y.views::REAL AS "value",
    'youtube.video.views' AS "metric"
FROM youtube_video_views AS y
WHERE
    y.id = ANY($1::VARCHAR[]) AND
    y.ts >= $2 AND
    y.ts <= $3
UNION ALL
SELECT
    y.id AS "id",
    y.title AS "title",
    y.ts AS "ts",
    y.likes::REAL AS "value",
    'youtube.video.likes' AS "metric"
FROM youtube_video_likes AS y
WHERE
    y.id = ANY($1::VARCHAR[]) AND
    y.ts >= $2 AND
    y.ts <= $3
UNION ALL
SELECT
    y.id AS "id",
    y.title AS "title",
    y.ts AS "ts",
    y.comments::REAL AS "value",
    'youtube.video.comments' AS "metric"
FROM youtube_video_comments AS y
WHERE
    y.id = ANY($1::VARCHAR[]) AND
    y.ts >= $2 AND
    y.ts <= $3
`

type GetYouTubeVideoMetricsByIDsParams struct {
	Ids     []string           `json:"ids"`
	TsStart pgtype.Timestamptz `json:"ts_start"`
	TsEnd   pgtype.Timestamptz `json:"ts_end"`
}

type GetYouTubeVideoMetricsByIDsRow struct {
	ID     string             `json:"id"`
	Title  string             `json:"title"`
	Ts     pgtype.Timestamptz `json:"ts"`
	Value  float32            `json:"value"`
	Metric string             `json:"metric"`
}

func (q *Queries) GetYouTubeVideoMetricsByIDs(ctx context.Context, arg GetYouTubeVideoMetricsByIDsParams) ([]GetYouTubeVideoMetricsByIDsRow, error) {
	rows, err := q.db.Query(ctx, getYouTubeVideoMetricsByIDs, arg.Ids, arg.TsStart, arg.TsEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetYouTubeVideoMetricsByIDsRow
	for rows.Next() {
		var i GetYouTubeVideoMetricsByIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Ts,
			&i.Value,
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

const getYouTubeVideoMetricsByIDsBucketed = `-- name: GetYouTubeVideoMetricsByIDsBucketed :many
SELECT
    y.id AS "id",
    FIRST(y.title, y.ts) AS "title",
    time_bucket(INTERVAL '60 min', y.ts) AS "ts",
    MAX(y.views::REAL) AS "value",
    'youtube.video.views' AS "metric"
FROM youtube_video_views AS y
GROUP BY y.id, y.ts
HAVING
    y.id = ANY($1::VARCHAR[]) AND
    y.ts >= $2 AND
    y.ts <= $3
`

type GetYouTubeVideoMetricsByIDsBucketedParams struct {
	Ids     []string           `json:"ids"`
	TsStart pgtype.Timestamptz `json:"ts_start"`
	TsEnd   pgtype.Timestamptz `json:"ts_end"`
}

type GetYouTubeVideoMetricsByIDsBucketedRow struct {
	ID     string      `json:"id"`
	Title  interface{} `json:"title"`
	Ts     interface{} `json:"ts"`
	Value  interface{} `json:"value"`
	Metric string      `json:"metric"`
}

func (q *Queries) GetYouTubeVideoMetricsByIDsBucketed(ctx context.Context, arg GetYouTubeVideoMetricsByIDsBucketedParams) ([]GetYouTubeVideoMetricsByIDsBucketedRow, error) {
	rows, err := q.db.Query(ctx, getYouTubeVideoMetricsByIDsBucketed, arg.Ids, arg.TsStart, arg.TsEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetYouTubeVideoMetricsByIDsBucketedRow
	for rows.Next() {
		var i GetYouTubeVideoMetricsByIDsBucketedRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Ts,
			&i.Value,
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
INSERT INTO youtube_video_comments (id, title, ts, comments)
VALUES ($1, $2, NOW()::TIMESTAMPTZ, $3)
`

type InsertYouTubeVideoCommentsParams struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Comments int32  `json:"comments"`
}

func (q *Queries) InsertYouTubeVideoComments(ctx context.Context, arg InsertYouTubeVideoCommentsParams) error {
	_, err := q.db.Exec(ctx, insertYouTubeVideoComments, arg.ID, arg.Title, arg.Comments)
	return err
}

const insertYouTubeVideoLikes = `-- name: InsertYouTubeVideoLikes :exec
INSERT INTO youtube_video_likes (id, title, ts, likes)
VALUES ($1, $2, NOW()::TIMESTAMPTZ, $3)
`

type InsertYouTubeVideoLikesParams struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Likes int32  `json:"likes"`
}

func (q *Queries) InsertYouTubeVideoLikes(ctx context.Context, arg InsertYouTubeVideoLikesParams) error {
	_, err := q.db.Exec(ctx, insertYouTubeVideoLikes, arg.ID, arg.Title, arg.Likes)
	return err
}

const insertYouTubeVideoViews = `-- name: InsertYouTubeVideoViews :exec
INSERT INTO youtube_video_views (id, title, ts, views)
VALUES ($1, $2, NOW()::TIMESTAMPTZ, $3)
`

type InsertYouTubeVideoViewsParams struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Views int32  `json:"views"`
}

func (q *Queries) InsertYouTubeVideoViews(ctx context.Context, arg InsertYouTubeVideoViewsParams) error {
	_, err := q.db.Exec(ctx, insertYouTubeVideoViews, arg.ID, arg.Title, arg.Views)
	return err
}
