-- name: GetRedditUserSubscriptions :many
SELECT name
FROM reddit_user_subscriptions;

-- name: GetRedditSubredditSubscriptions :many
SELECT name
FROM reddit_subreddit_subscriptions;