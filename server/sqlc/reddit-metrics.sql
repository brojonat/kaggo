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

-- name: GetRedditPostMetricsByID :many
SELECT *, 'reddit.post.score' AS "metric"
FROM reddit_post_score rps
WHERE rps.id = @id
UNION ALL
SELECT *, 'reddit.post.ratio' AS "metric"
FROM reddit_post_ratio rpr
WHERE rpr.id = @id;

-- name: GetRedditCommentMetricsByID :many
SELECT *, 'reddit.comment.score' AS "metric"
FROM reddit_comment_score rcs
WHERE rcs.id = @id
UNION ALL
SELECT *, 'reddit.comment.controversiality' AS "metric"
FROM reddit_comment_controversiality rcc
WHERE rcc.id = @id;
