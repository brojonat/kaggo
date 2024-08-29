package jsonb

type MetadataJSON struct {
	// these fields should always be set
	ID   string   `json:"id"`
	Link string   `json:"link"`
	Tags []string `json:"tags"`
	// the remaining fields may or may not be set in the JSON read from the DB,
	// and may or may not be present in the JSON written to a client.
	Owner       string `json:"owner,omitempty"`
	Title       string `json:"title,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Subreddit   string `json:"subreddit,omitempty"`
	GameID      string `json:"game_id,omitempty"`
	Broadcaster string `json:"broadcaster,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UserMetadataJSON struct{}
