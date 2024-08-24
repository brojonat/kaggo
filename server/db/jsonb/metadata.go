package jsonb

type MetadataJSON struct {
	// these fields should always be set
	ID   string `json:"id"`
	Link string `json:"link"`
	// the remaining fields may or may not be set in the JSON read from the DB,
	// and may or may not be present in the JSON written to a client.
	Owner   string   `json:"owner,omitempty"`
	Author  string   `json:"author,omitempty"`
	Title   string   `json:"title,omitempty"`
	Comment string   `json:"comment,omitempty"`
	Tags    []string `json:"tags"`
	// author? subreddit? created at?
}
