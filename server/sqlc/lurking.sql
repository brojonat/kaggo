-- name: InsertYouTubeChannelSubscription :exec
INSERT INTO youtube_channel_subscriptions (id)
VALUES (@id);

-- name: GetYouTubeChannelSubscriptions :many
SELECT id
FROM youtube_channel_subscriptions;

-- name: YouTubeChannelSubscriptionExists :one
SELECT 1 AS "exists"
FROM youtube_channel_subscriptions
WHERE id = @id;
