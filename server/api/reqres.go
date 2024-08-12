package api

type DefaultJSONResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type InternalMetricPayload struct {
	ID    string `json:"id"`
	Value int    `json:"value"`
}

type YouTubeVideoMetricPayload struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	SetViews    bool   `json:"set_views"`
	Views       int    `json:"views"`
	SetComments bool   `json:"set_comments"`
	Comments    int    `json:"comments"`
	SetLikes    bool   `json:"set_likes"`
	Likes       int    `json:"likes"`
}

type KaggleMetricPayload struct {
	Slug         string `json:"slug"`
	SetVotes     bool   `json:"set_votes"`
	Votes        int    `json:"votes"`
	SetDownloads bool   `json:"set_downloads"`
	Downloads    int    `json:"downloads"`
}
