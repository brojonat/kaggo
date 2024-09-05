// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: lurking.sql

package dbgen

import (
	"context"
)

const getRedditSubredditSubscriptions = `-- name: GetRedditSubredditSubscriptions :many
SELECT name
FROM reddit_subreddit_subscriptions
`

func (q *Queries) GetRedditSubredditSubscriptions(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getRedditSubredditSubscriptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRedditUserSubscriptions = `-- name: GetRedditUserSubscriptions :many
SELECT name
FROM reddit_user_subscriptions
`

func (q *Queries) GetRedditUserSubscriptions(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getRedditUserSubscriptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getYouTubeChannelSubscriptions = `-- name: GetYouTubeChannelSubscriptions :many
SELECT id
FROM youtube_channel_subscriptions
`

func (q *Queries) GetYouTubeChannelSubscriptions(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getYouTubeChannelSubscriptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertRedditUserSubscription = `-- name: InsertRedditUserSubscription :exec
INSERT INTO reddit_user_subscriptions (name)
VALUES ($1)
`

func (q *Queries) InsertRedditUserSubscription(ctx context.Context, name string) error {
	_, err := q.db.Exec(ctx, insertRedditUserSubscription, name)
	return err
}

const insertYouTubeChannelSubscription = `-- name: InsertYouTubeChannelSubscription :exec
INSERT INTO youtube_channel_subscriptions (id)
VALUES ($1)
`

func (q *Queries) InsertYouTubeChannelSubscription(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, insertYouTubeChannelSubscription, id)
	return err
}

const youTubeChannelSubscriptionExists = `-- name: YouTubeChannelSubscriptionExists :one
SELECT 1 AS "exists"
FROM youtube_channel_subscriptions
WHERE id = $1
`

func (q *Queries) YouTubeChannelSubscriptionExists(ctx context.Context, id string) (int32, error) {
	row := q.db.QueryRow(ctx, youTubeChannelSubscriptionExists, id)
	var exists int32
	err := row.Scan(&exists)
	return exists, err
}
