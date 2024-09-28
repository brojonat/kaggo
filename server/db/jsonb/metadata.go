package jsonb

import "time"

type MetadataJSON struct {
	// these fields should always be set
	ID         string   `json:"id"`
	HumanLabel string   `json:"human_label"`
	Link       string   `json:"link"`
	Tags       []string `json:"tags"`
	// the remaining fields may or may not be set in the JSON read from the DB,
	// and may or may not be present in the JSON written to a client.
	Owner              string    `json:"owner,omitempty"`
	Title              string    `json:"title,omitempty"`
	Comment            string    `json:"comment,omitempty"`
	TSCreated          time.Time `json:"ts_created,omitempty"`
	UserID             string    `json:"user_id,omitempty"`
	ParentUserID       string    `json:"parent_user_id,omitempty"`
	ParentUserName     string    `json:"parent_user_name,omitempty"`
	ParentPostID       string    `json:"parent_post_id,omitempty"`
	ParentPostTitle    string    `json:"parent_post_title,omitempty"`
	ParentSubreddit    string    `json:"parent_subreddit,omitempty"`
	ParentChannelID    string    `json:"parent_channel_id,omitempty"`
	ParentChannelTitle string    `json:"parent_channel_title,omitempty"`
	GameID             string    `json:"game_id,omitempty"`
	Broadcaster        string    `json:"broadcaster,omitempty"`
	Duration           int       `json:"duration,omitempty"`
	DisplayName        string    `json:"display_name,omitempty"`
	Description        string    `json:"description,omitempty"`
}

type UserMetadataJSON struct{}
