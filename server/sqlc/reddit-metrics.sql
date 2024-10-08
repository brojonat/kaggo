-- name: InsertRedditPostScore :exec
INSERT INTO reddit_post_score (id, ts, score)
VALUES (@id, NOW()::TIMESTAMPTZ, @score);

-- name: InsertRedditPostRatio :exec
INSERT INTO reddit_post_ratio (id, ts, ratio)
VALUES (@id, NOW()::TIMESTAMPTZ, @ratio);

-- name: InsertRedditCommentScore :exec
INSERT INTO reddit_comment_score (id, ts, score)
VALUES (@id, NOW()::TIMESTAMPTZ, @score);

-- name: InsertRedditCommentControversiality :exec
INSERT INTO reddit_comment_controversiality (id, ts, controversiality)
VALUES (@id, NOW()::TIMESTAMPTZ, @controversiality);

-- name: InsertRedditSubredditSubscribers :exec
INSERT INTO reddit_subreddit_subscribers (id, ts, subscribers)
VALUES (@id, NOW()::TIMESTAMPTZ, @subscribers);

-- name: InsertRedditSubredditActiveUserCount :exec
INSERT INTO reddit_subreddit_active_user_count (id, ts, active_user_count)
VALUES (@id, NOW()::TIMESTAMPTZ, @active_user_count);

-- name: InsertRedditUserAwardeeKarma :exec
INSERT INTO reddit_user_awardee_karma (id, ts, karma)
VALUES (@id, NOW()::TIMESTAMPTZ, @karma);

-- name: InsertRedditUserAwarderKarma :exec
INSERT INTO reddit_user_awarder_karma (id, ts, karma)
VALUES (@id, NOW()::TIMESTAMPTZ, @karma);

-- name: InsertRedditUserCommentKarma :exec
INSERT INTO reddit_user_comment_karma (id, ts, karma)
VALUES (@id, NOW()::TIMESTAMPTZ, @karma);

-- name: InsertRedditUserLinkKarma :exec
INSERT INTO reddit_user_link_karma (id, ts, karma)
VALUES (@id, NOW()::TIMESTAMPTZ, @karma);

-- name: InsertRedditUserTotalKarma :exec
INSERT INTO reddit_user_total_karma (id, ts, karma)
VALUES (@id, NOW()::TIMESTAMPTZ, @karma);

-- name: GetRedditPostMetricsByIDs :many
SELECT
    id AS "id",
    ts AS "ts",
    score::REAL AS "value",
    'reddit.post.score' AS "metric"
FROM reddit_post_score AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end
UNION ALL
SELECT
    id AS "id",
    ts AS "ts",
    ratio::REAL AS "value",
    'reddit.post.ratio' AS "metric"
FROM reddit_post_ratio AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end;

-- name: GetRedditCommentMetricsByIDs :many
SELECT
    r.id AS "id",
    r.ts AS "ts",
    r.score::REAL AS "value",
    'reddit.comment.score' AS "metric"
FROM reddit_comment_score AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end
UNION ALL
SELECT
    r.id AS "id",
    r.ts AS "ts",
    r.controversiality::REAL AS "value",
    'reddit.comment.controversiality' AS "metric"
FROM reddit_comment_controversiality AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end;

-- name: GetRedditSubredditMetricsByIDs :many
SELECT
    r.id AS "id",
    r.ts AS "ts",
    r.subscribers::REAL AS "value",
    'reddit.subreddit.subscribers' AS "metric"
FROM reddit_subreddit_subscribers AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end
UNION ALL
SELECT
    r.id AS "id",
    r.ts AS "ts",
    r.active_user_count::REAL AS "value",
    'reddit.subreddit.active-user-count' AS "metric"
FROM reddit_subreddit_active_user_count AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end;

-- name: GetRedditUserMetricsByIDs :many
SELECT
    r.id AS "id",
    r.ts AS "ts",
    r.karma::REAL AS "value",
    'reddit.user.awardee-karma' AS "metric"
FROM reddit_user_awardee_karma AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end
UNION ALL
SELECT
    r.id AS "id",
    r.ts AS "ts",
    r.karma::REAL AS "value",
    'reddit.user.awarder-karma' AS "metric"
FROM reddit_user_awarder_karma AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end
UNION ALL
SELECT
    r.id AS "id",
    r.ts AS "ts",
    r.karma::REAL AS "value",
    'reddit.user.comment-karma' AS "metric"
FROM reddit_user_comment_karma AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end
UNION ALL
SELECT
    r.id AS "id",
    r.ts AS "ts",
    r.karma::REAL AS "value",
    'reddit.user.link-karma' AS "metric"
FROM reddit_user_link_karma AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end
UNION ALL
SELECT
    r.id AS "id",
    r.ts AS "ts",
    r.karma::REAL AS "value",
    'reddit.user.total-karma' AS "metric"
FROM reddit_user_total_karma AS r
WHERE
    r.id ILIKE ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end;

-- Reddit Post Bucketed Metrics

-- name: GetRedditPostMetricsByIDsBucket15Min :many
SELECT *, 'reddit.post.score' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(score::REAL) AS "value"
	FROM reddit_post_score
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.post.ratio' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(ratio::REAL) AS "value"
	FROM reddit_post_ratio
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditPostMetricsByIDsBucket1Hr :many
SELECT *, 'reddit.post.score' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(score::REAL) AS value
	FROM reddit_post_score
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.post.ratio' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(ratio::REAL) AS value
	FROM reddit_post_ratio
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditPostMetricsByIDsBucket8Hr :many
SELECT *, 'reddit.post.score' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS bucket,
	    MAX(score::REAL) AS value
	FROM reddit_post_score
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.post.ratio' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(ratio::REAL) AS "value"
	FROM reddit_post_ratio
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditPostMetricsByIDsBucket1Day :many
SELECT *, 'reddit.post.score' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(score::REAL) AS value
	FROM reddit_post_score
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.post.ratio' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(ratio::REAL) AS value
	FROM reddit_post_ratio
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;


-- Reddit Comment Bucketed Metrics


-- name: GetRedditCommentMetricsByIDsBucket15Min :many
SELECT *, 'reddit.comment.score' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(score::REAL) AS "value"
	FROM reddit_comment_score
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.comment.controversiality' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(controversiality::REAL) AS "value"
	FROM reddit_comment_controversiality
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditCommentMetricsByIDsBucket1Hr :many
SELECT *, 'reddit.comment.score' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(score::REAL) AS value
	FROM reddit_comment_score
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.comment.controversiality' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(controversiality::REAL) AS value
	FROM reddit_comment_controversiality
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditCommentMetricsByIDsBucket8Hr :many
SELECT *, 'reddit.comment.score' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS bucket,
	    MAX(score::REAL) AS value
	FROM reddit_comment_score
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.comment.controversiality' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(controversiality::REAL) AS "value"
	FROM reddit_comment_controversiality
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditCommentMetricsByIDsBucket1Day :many
SELECT *, 'reddit.comment.score' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(score::REAL) AS value
	FROM reddit_comment_score
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.comment.controversiality' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(controversiality::REAL) AS value
	FROM reddit_comment_controversiality
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;


-- Reddit Subreddit Bucketed Metrics


-- name: GetRedditSubredditMetricsByIDsBucket15Min :many
SELECT *, 'reddit.subreddit.subscribers' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(subscribers::REAL) AS "value"
	FROM reddit_subreddit_subscribers
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.subreddit.active-user-count' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(active_user_count::REAL) AS "value"
	FROM reddit_subreddit_active_user_count
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditSubredditMetricsByIDsBucket1Hr :many
SELECT *, 'reddit.subreddit.subscribers' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(subscribers::REAL) AS value
	FROM reddit_subreddit_subscribers
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.subreddit.active-user-count' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(active_user_count::REAL) AS value
	FROM reddit_subreddit_active_user_count
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditSubredditMetricsByIDsBucket8Hr :many
SELECT *, 'reddit.subreddit.subscribers' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS bucket,
	    MAX(subscribers::REAL) AS value
	FROM reddit_subreddit_subscribers
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.subreddit.active-user-count' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(active_user_count::REAL) AS "value"
	FROM reddit_subreddit_active_user_count
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditSubredditMetricsByIDsBucket1Day :many
SELECT *, 'reddit.subreddit.subscribers' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(subscribers::REAL) AS value
	FROM reddit_subreddit_subscribers
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.subreddit.active-user-count' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(active_user_count::REAL) AS value
	FROM reddit_subreddit_active_user_count
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;


-- Reddit User Bucketed Metrics


-- name: GetRedditUserMetricsByIDsBucket15Min :many
SELECT *, 'reddit.user.awardee-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_awardee_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.awarder-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_awarder_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.comment-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_comment_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.link-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_link_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.total-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_total_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditUserMetricsByIDsBucket1Hr :many
SELECT *, 'reddit.user.awardee-karma' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(karma::REAL) AS value
	FROM reddit_user_awardee_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.awarder-karma' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(karma::REAL) AS value
	FROM reddit_user_awarder_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.comment-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_comment_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.link-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_link_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.total-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_total_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditUserMetricsByIDsBucket8Hr :many
SELECT *, 'reddit.user.awardee-karma' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS bucket,
	    MAX(karma::REAL) AS value
	FROM reddit_user_awardee_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.awarder-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_awarder_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.comment-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_comment_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.link-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_link_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.total-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_total_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetRedditUserMetricsByIDsBucket1Day :many
SELECT *, 'reddit.user.awardee-karma' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(karma::REAL) AS value
	FROM reddit_user_awardee_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.awarder-karma' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(karma::REAL) AS value
	FROM reddit_user_awarder_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
	UNION ALL
SELECT *, 'reddit.user.comment-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_comment_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.link-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_link_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'reddit.user.total-karma' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS "bucket",
	    MAX(karma::REAL) AS "value"
	FROM reddit_user_total_karma
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id ILIKE ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;
