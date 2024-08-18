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
	Title       string `json:"title"`
	SetViews    bool   `json:"set_views"`
	Views       int    `json:"views"`
	SetComments bool   `json:"set_comments"`
	Comments    int    `json:"comments"`
	SetLikes    bool   `json:"set_likes"`
	Likes       int    `json:"likes"`
}

type KaggleNotebookMetricPayload struct {
	ID           string `json:"id"`
	SetViews     bool   `json:"set_views"`
	Views        int    `json:"views"`
	SetVotes     bool   `json:"set_votes"`
	Votes        int    `json:"votes"`
	SetDownloads bool   `json:"set_downloads"`
	Downloads    int    `json:"downloads"`
}

type KaggleDatasetMetricPayload struct {
	ID           string `json:"id"`
	SetViews     bool   `json:"set_views"`
	Views        int    `json:"views"`
	SetVotes     bool   `json:"set_votes"`
	Votes        int    `json:"votes"`
	SetDownloads bool   `json:"set_downloads"`
	Downloads    int    `json:"downloads"`
}

type RedditPostMetricPayload struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	SetScore bool    `json:"set_score"`
	Score    int     `json:"score"`
	SetRatio bool    `json:"set_ratio"`
	Ratio    float32 `json:"ratio"`
}

type RedditCommentMetricPayload struct {
	ID                  string  `json:"id"`
	SetScore            bool    `json:"set_score"`
	Score               int     `json:"score"`
	SetControversiality bool    `json:"set_controversiality"`
	Controversiality    float32 `json:"controversiality"`
}
