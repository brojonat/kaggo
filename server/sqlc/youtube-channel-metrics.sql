-- name: InsertYouTubeChannelViews :exec
INSERT INTO youtube_channel_views (id, ts, views)
VALUES (@id, NOW()::TIMESTAMPTZ, @views);

-- name: InsertYouTubeChannelSubscribers :exec
INSERT INTO youtube_channel_subscribers (id, ts, subscribers)
VALUES (@id, NOW()::TIMESTAMPTZ, @subscribers);

-- name: InsertYouTubeChannelVideos :exec
INSERT INTO youtube_channel_videos (id, ts, videos)
VALUES (@id, NOW()::TIMESTAMPTZ, @videos);

-- name: GetYouTubeChannelMetricsByIDs :many
SELECT
    y.id AS "id",
    y.ts AS "ts",
    y.views::REAL AS "value",
    'youtube.channel.views' AS "metric"
FROM youtube_channel_views AS y
WHERE
    y.id = ANY(@ids::VARCHAR[]) AND
    y.ts >= @ts_start AND
    y.ts <= @ts_end
UNION ALL
SELECT
    y.id AS "id",
    y.ts AS "ts",
    y.subscribers::REAL AS "value",
    'youtube.channel.subscribers' AS "metric"
FROM youtube_channel_subscribers AS y
WHERE
    y.id = ANY(@ids::VARCHAR[]) AND
    y.ts >= @ts_start AND
    y.ts <= @ts_end
UNION ALL
SELECT
    y.id AS "id",
    y.ts AS "ts",
    y.videos::REAL AS "value",
    'youtube.channel.videos' AS "metric"
FROM youtube_channel_videos AS y
WHERE
    y.id = ANY(@ids::VARCHAR[]) AND
    y.ts >= @ts_start AND
    y.ts <= @ts_end;

-- name: GetYouTubeChannelMetricsByIDsBucketed :many
SELECT
    y.id AS "id",
    FIRST(y.title, y.ts) AS "title",
    time_bucket(INTERVAL '60 min', y.ts) AS "ts",
    MAX(y.views::REAL) AS "value",
    'youtube.channel.views' AS "metric"
FROM youtube_channel_views AS y
GROUP BY y.id, y.ts
HAVING
    y.id = ANY(@ids::VARCHAR[]) AND
    y.ts >= @ts_start AND
    y.ts <= @ts_end;

-- name: GetYouTubeChannelMetricsByIDsBucket15Min :many
SELECT *, 'youtube.channel.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM youtube_channel_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.channel.subscribers' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(subscribers::REAL) AS "value"
	FROM youtube_channel_subscribers
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.channel.videos' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(videos::REAL) AS "value"
	FROM youtube_channel_videos
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetYouTubeChannelMetricsByIDsBucket1Hr :many
SELECT *, 'youtube.channel.views' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(views::REAL) AS value
	FROM youtube_channel_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.channel.subscribers' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(subscribers::REAL) AS value
	FROM youtube_channel_subscribers
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.channel.videos' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(videos::REAL) AS value
	FROM youtube_channel_videos
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetYouTubeChannelMetricsByIDsBucket8Hr :many
SELECT *, 'youtube.channel.views' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS bucket,
	    MAX(views::REAL) AS value
	FROM youtube_channel_views AS yvv
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.channel.subscribers' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(subscribers::REAL) AS "value"
	FROM youtube_channel_subscribers
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.channel.videos' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(videos::REAL) AS "value"
	FROM youtube_channel_videos
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetYouTubeChannelMetricsByIDsBucket1Day :many
SELECT *, 'youtube.channel.views' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(views::REAL) AS value
	FROM youtube_channel_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.channel.subscribers' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(subscribers::REAL) AS value
	FROM youtube_channel_subscribers
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'youtube.channel.videos' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(videos::REAL) AS value
	FROM youtube_channel_videos
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;