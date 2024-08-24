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