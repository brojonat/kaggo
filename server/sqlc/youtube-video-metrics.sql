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
    y.id ILIKE ANY(@ids::VARCHAR[]) AND
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
    y.id ILIKE ANY(@ids::VARCHAR[]) AND
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
    y.id ILIKE ANY(@ids::VARCHAR[]) AND
    y.ts >= @ts_start AND
    y.ts <= @ts_end;

-- name: GetYouTubeVideoMetricsByIDsBucket15Min :many
SELECT *, 'youtube.video.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM youtube_video_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.video.likes' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(likes::REAL) AS "value"
	FROM youtube_video_likes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.video.comments' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(comments::REAL) AS "value"
	FROM youtube_video_comments
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetYouTubeVideoMetricsByIDsBucket1Hr :many
SELECT *, 'youtube.video.views' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(views::REAL) AS value
	FROM youtube_video_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.video.likes' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(likes::REAL) AS value
	FROM youtube_video_likes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.video.comments' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(comments::REAL) AS value
	FROM youtube_video_comments
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetYouTubeVideoMetricsByIDsBucket8Hr :many
SELECT *, 'youtube.video.views' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS bucket,
	    MAX(views::REAL) AS value
	FROM youtube_video_views AS yvv
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.video.likes' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(likes::REAL) AS "value"
	FROM youtube_video_likes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.video.comments' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(comments::REAL) AS "value"
	FROM youtube_video_comments
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetYouTubeVideoMetricsByIDsBucket1Day :many
SELECT *, 'youtube.video.views' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(views::REAL) AS value
	FROM youtube_video_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.video.likes' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(likes::REAL) AS value
	FROM youtube_video_likes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.video.comments' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(comments::REAL) AS value
	FROM youtube_video_comments
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;
