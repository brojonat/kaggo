-- name: InsertYouTubeVideoViews :exec
INSERT INTO youtube_video_views (id, title, ts, views)
VALUES (@id, @title, NOW()::TIMESTAMPTZ, @views);

-- name: InsertYouTubeVideoComments :exec
INSERT INTO youtube_video_comments (id, title, ts, comments)
VALUES (@id, @title, NOW()::TIMESTAMPTZ, @comments);

-- name: InsertYouTubeVideoLikes :exec
INSERT INTO youtube_video_likes (id, title, ts, likes)
VALUES (@id, @title, NOW()::TIMESTAMPTZ, @likes);

-- name: GetYouTubeVideoMetricsByIDs :many
SELECT
    y.id AS "id",
    y.title AS "title",
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
    y.title AS "title",
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
    y.title AS "title",
    y.ts AS "ts",
    y.comments::REAL AS "value",
    'youtube.video.comments' AS "metric"
FROM youtube_video_comments AS y
WHERE
    y.id = ANY(@ids::VARCHAR[]) AND
    y.ts >= @ts_start AND
    y.ts <= @ts_end;
