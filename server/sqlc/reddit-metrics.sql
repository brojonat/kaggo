-- name: InsertRedditPostScore :exec
INSERT INTO reddit_post_score (id, title, ts, score)
VALUES (@id, @title, NOW()::TIMESTAMPTZ, @score);

-- name: InsertRedditPostRatio :exec
INSERT INTO reddit_post_ratio (id, title, ts, ratio)
VALUES (@id, @title, NOW()::TIMESTAMPTZ, @ratio);

-- name: InsertRedditCommentScore :exec
INSERT INTO reddit_comment_score (id, ts, score)
VALUES (@id, NOW()::TIMESTAMPTZ, @score);

-- name: InsertRedditCommentControversiality :exec
INSERT INTO reddit_comment_controversiality (id, ts, controversiality)
VALUES (@id, NOW()::TIMESTAMPTZ, @controversiality);

-- name: GetRedditPostMetricsByIDs :many
SELECT
    id AS "id",
    title AS "title",
    ts AS "ts",
    score::REAL AS "value",
    'reddit.post.score' AS "metric"
FROM reddit_post_score AS r
WHERE
    r.id = ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end
UNION ALL
SELECT
    id AS "id",
    title AS "title",
    ts AS "ts",
    ratio::REAL AS "value",
    'reddit.post.ratio' AS "metric"
FROM reddit_post_ratio AS r
WHERE
    r.id = ANY(@ids::VARCHAR[]) AND
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
    r.id = ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end
UNION ALL
SELECT
    r.id AS "id",
    r.ts AS "ts",
    r.controversiality::REAL AS "controversiality",
    'reddit.comment.controversiality' AS "metric"
FROM reddit_comment_controversiality AS r
WHERE
    r.id = ANY(@ids::VARCHAR[]) AND
    r.ts >= @ts_start AND
    r.ts <= @ts_end;

