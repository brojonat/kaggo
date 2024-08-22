-- name: InsertYouTubeVideoViews :exec
INSERT INTO youtube_video_views (id, ts, views)
VALUES (@id, NOW()::TIMESTAMPTZ, @views);

-- name: InsertYouTubeVideoComments :exec
INSERT INTO youtube_video_comments (id, ts, comments)
VALUES (@id, NOW()::TIMESTAMPTZ, @comments);

-- name: InsertYouTubeVideoLikes :exec
INSERT INTO youtube_video_likes (id, ts, likes)
VALUES (@id, NOW()::TIMESTAMPTZ, @likes);

-- name: GetYouTubeVideoMetricsByIDs :many
SELECT
    y.id AS "id",
    y.ts AS "ts",
    y.views::REAL AS "value",
    'youtube.video.views' AS "metric"
FROM youtube_video_views AS y
WHERE
    y.id = ANY(@ids::VARCHAR[]) AND
    y.ts >= @ts_start AND
    y.ts <= @ts_end
UNION ALL
SELECT
    y.id AS "id",
    y.ts AS "ts",
    y.likes::REAL AS "value",
    'youtube.video.likes' AS "metric"
FROM youtube_video_likes AS y
WHERE
    y.id = ANY(@ids::VARCHAR[]) AND
    y.ts >= @ts_start AND
    y.ts <= @ts_end
UNION ALL
SELECT
    y.id AS "id",
    y.ts AS "ts",
    y.comments::REAL AS "value",
    'youtube.video.comments' AS "metric"
FROM youtube_video_comments AS y
WHERE
    y.id = ANY(@ids::VARCHAR[]) AND
    y.ts >= @ts_start AND
    y.ts <= @ts_end;

-- name: GetYouTubeVideoMetricsByIDsBucketed :many
SELECT
    y.id AS "id",
    FIRST(y.title, y.ts) AS "title",
    time_bucket(INTERVAL '60 min', y.ts) AS "ts",
    MAX(y.views::REAL) AS "value",
    'youtube.video.views' AS "metric"
FROM youtube_video_views AS y
GROUP BY y.id, y.ts
HAVING
    y.id = ANY(@ids::VARCHAR[]) AND
    y.ts >= @ts_start AND
    y.ts <= @ts_end;