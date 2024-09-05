-- name: InsertYouTubeChannelSubscription :exec
INSERT INTO youtube_channel_subscriptions (id)
VALUES (@id);

-- name: InsertRedditUserSubscription :exec
INSERT INTO reddit_user_subscriptions (name)
VALUES (@name);

-- name: GetYouTubeChannelSubscriptions :many
SELECT id
FROM youtube_channel_subscriptions;

-- name: YouTubeChannelSubscriptionExists :one
SELECT 1 AS "exists"
FROM youtube_channel_subscriptions
WHERE id = @id;

-- name: GetRedditUserSubscriptions :many
SELECT name
FROM reddit_user_subscriptions;

-- name: GetRedditSubredditSubscriptions :many
SELECT name
FROM reddit_subreddit_subscriptions;
