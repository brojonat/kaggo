-- name: InsertTwitchClipViews :exec
INSERT INTO twitch_clip_views (id, ts, views)
VALUES (@id, NOW()::TIMESTAMPTZ, @views);

-- name: InsertTwitchVideoViews :exec
INSERT INTO twitch_video_views (id, ts, views)
VALUES (@id, NOW()::TIMESTAMPTZ, @views);

-- name: InsertTwitchStreamViews :exec
INSERT INTO twitch_stream_views (id, ts, views)
VALUES (@id, NOW()::TIMESTAMPTZ, @views);

-- name: InsertTwitchUserPastDecAvgViews :exec
INSERT INTO twitch_user_past_dec_avg_views (id, ts, avg_views)
VALUES (@id, NOW()::TIMESTAMPTZ, @avg_views);

-- name: InsertTwitchUserPastDecMedViews :exec
INSERT INTO twitch_user_past_dec_med_views (id, ts, med_views)
VALUES (@id, NOW()::TIMESTAMPTZ, @med_views);

-- name: InsertTwitchUserPastDecStdViews :exec
INSERT INTO twitch_user_past_dec_std_views (id, ts, std_views)
VALUES (@id, NOW()::TIMESTAMPTZ, @std_views);

-- name: InsertTwitchUserPastDecAvgDuration :exec
INSERT INTO twitch_user_past_dec_avg_duration (id, ts, avg_duration)
VALUES (@id, NOW()::TIMESTAMPTZ, @avg_duration);

-- name: InsertTwitchUserPastDecMedDuration :exec
INSERT INTO twitch_user_past_dec_med_duration (id, ts, med_duration)
VALUES (@id, NOW()::TIMESTAMPTZ, @med_duration);

-- name: InsertTwitchUserPastDecStdDuration :exec
INSERT INTO twitch_user_past_dec_std_duration (id, ts, std_duration)
VALUES (@id, NOW()::TIMESTAMPTZ, @std_duration);


-- name: GetTwitchClipMetricsByIDsBucket15Min :many
SELECT *, 'twitch.clip.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_clip_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchClipMetricsByIDsBucket1Hr :many
SELECT *, 'twitch.clip.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_clip_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchClipMetricsByIDsBucket8Hr :many
SELECT *, 'twitch.clip.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_clip_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchClipMetricsByIDsBucket1Day :many
SELECT *, 'twitch.clip.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_clip_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchVideoMetricsByIDsBucket15Min :many
SELECT *, 'twitch.video.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_video_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchVideoMetricsByIDsBucket1Hr :many
SELECT *, 'twitch.video.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_video_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchVideoMetricsByIDsBucket8Hr :many
SELECT *, 'twitch.video.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_video_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchVideoMetricsByIDsBucket1Day :many
SELECT *, 'twitch.video.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_video_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchStreamMetricsByIDsBucket15Min :many
SELECT *, 'twitch.stream.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_stream_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchStreamMetricsByIDsBucket1Hr :many
SELECT *, 'twitch.stream.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_stream_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchStreamMetricsByIDsBucket8Hr :many
SELECT *, 'twitch.stream.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_stream_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchStreamMetricsByIDsBucket1Day :many
SELECT *, 'twitch.stream.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM twitch_stream_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchUserPastDecMetricsByIDsBucket15Min :many
SELECT *, 'twitch.user-past-dec.avg-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(avg_views::REAL) AS "value"
	FROM twitch_user_past_dec_avg_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.med-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(med_views::REAL) AS "value"
	FROM twitch_user_past_dec_med_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.std-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(std_views::REAL) AS "value"
	FROM twitch_user_past_dec_std_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.avg-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(avg_duration::REAL) AS "value"
	FROM twitch_user_past_dec_avg_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.med-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(med_duration::REAL) AS "value"
	FROM twitch_user_past_dec_med_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.std-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(std_duration::REAL) AS "value"
	FROM twitch_user_past_dec_std_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;


-- name: GetTwitchUserPastDecMetricsByIDsBucket1Hr :many
SELECT *, 'twitch.user-past-dec.avg-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(avg_views::REAL) AS "value"
	FROM twitch_user_past_dec_avg_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.med-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(med_views::REAL) AS "value"
	FROM twitch_user_past_dec_med_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.std-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(std_views::REAL) AS "value"
	FROM twitch_user_past_dec_std_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.avg-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(avg_duration::REAL) AS "value"
	FROM twitch_user_past_dec_avg_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.med-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(med_duration::REAL) AS "value"
	FROM twitch_user_past_dec_med_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.std-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(std_duration::REAL) AS "value"
	FROM twitch_user_past_dec_std_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;


-- name: GetTwitchUserPastDecMetricsByIDsBucket8Hr :many
SELECT *, 'twitch.user-past-dec.avg-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(avg_views::REAL) AS "value"
	FROM twitch_user_past_dec_avg_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.med-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(med_views::REAL) AS "value"
	FROM twitch_user_past_dec_med_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.std-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(std_views::REAL) AS "value"
	FROM twitch_user_past_dec_std_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.avg-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(avg_duration::REAL) AS "value"
	FROM twitch_user_past_dec_avg_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.med-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(med_duration::REAL) AS "value"
	FROM twitch_user_past_dec_med_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.std-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(std_duration::REAL) AS "value"
	FROM twitch_user_past_dec_std_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetTwitchUserPastDecMetricsByIDsBucket1Day :many
SELECT *, 'twitch.user-past-dec.avg-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(avg_views::REAL) AS "value"
	FROM twitch_user_past_dec_avg_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.med-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(med_views::REAL) AS "value"
	FROM twitch_user_past_dec_med_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.std-views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(std_views::REAL) AS "value"
	FROM twitch_user_past_dec_std_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.avg-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(avg_duration::REAL) AS "value"
	FROM twitch_user_past_dec_avg_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.med-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(med_duration::REAL) AS "value"
	FROM twitch_user_past_dec_med_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'twitch.user-past-dec.std-duration' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(std_duration::REAL) AS "value"
	FROM twitch_user_past_dec_std_duration
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;
