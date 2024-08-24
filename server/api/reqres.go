package api

import (
	"github.com/brojonat/kaggo/server/db/jsonb"
	"go.temporal.io/sdk/client"
)

type DefaultJSONResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type GenericScheduleRequestPayload struct {
	RequestKind string              `json:"request_kind"`
	ID          string              `json:"id"`
	Schedule    client.ScheduleSpec `json:"schedule_spec,omitempty"`
}

// Contains arbitrary data that metric (or metadata) handlers want to pass back to the server.
type MetricQueryInternalData struct {
	XRatelimitUsed      string `json:"x_ratelimit_used"`
	XRatelimitRemaining string `json:"x_ratelimit_remaining"`
	XRatelimitReset     string `json:"x_ratelimit_reset"`
}

type MetricMetadataPayload struct {
	ID           string                  `json:"id"`
	RequestKind  string                  `json:"request_kind"`
	Data         jsonb.MetadataJSON      `json:"data"`
	InternalData MetricQueryInternalData `json:"internal_data"`
}

type InternalMetricPayload struct {
	ID           string                  `json:"id"`
	Value        int                     `json:"value"`
	InternalData MetricQueryInternalData `json:"internal_data"`
}

type YouTubeVideoMetricPayload struct {
	ID           string                  `json:"id"`
	SetViews     bool                    `json:"set_views"`
	Views        int                     `json:"views"`
	SetComments  bool                    `json:"set_comments"`
	Comments     int                     `json:"comments"`
	SetLikes     bool                    `json:"set_likes"`
	Likes        int                     `json:"likes"`
	InternalData MetricQueryInternalData `json:"internal_data"`
}

type YouTubeChannelMetricPayload struct {
	ID             string                  `json:"id"`
	SetViews       bool                    `json:"set_views"`
	Views          int                     `json:"views"`
	SetSubscribers bool                    `json:"set_subscribers"`
	Subscribers    int                     `json:"subscribers"`
	SetVideos      bool                    `json:"set_videos"`
	Videos         int                     `json:"videos"`
	InternalData   MetricQueryInternalData `json:"internal_data"`
}

type KaggleNotebookMetricPayload struct {
	ID           string                  `json:"id"`
	SetViews     bool                    `json:"set_views"`
	Views        int                     `json:"views"`
	SetVotes     bool                    `json:"set_votes"`
	Votes        int                     `json:"votes"`
	SetDownloads bool                    `json:"set_downloads"`
	Downloads    int                     `json:"downloads"`
	InternalData MetricQueryInternalData `json:"internal_data"`
}

type KaggleDatasetMetricPayload struct {
	ID           string                  `json:"id"`
	SetViews     bool                    `json:"set_views"`
	Views        int                     `json:"views"`
	SetVotes     bool                    `json:"set_votes"`
	Votes        int                     `json:"votes"`
	SetDownloads bool                    `json:"set_downloads"`
	Downloads    int                     `json:"downloads"`
	InternalData MetricQueryInternalData `json:"internal_data"`
}

type RedditPostMetricPayload struct {
	ID           string                  `json:"id"`
	SetScore     bool                    `json:"set_score"`
	Score        int                     `json:"score"`
	SetRatio     bool                    `json:"set_ratio"`
	Ratio        float32                 `json:"ratio"`
	InternalData MetricQueryInternalData `json:"internal_data"`
}

type RedditCommentMetricPayload struct {
	ID                  string                  `json:"id"`
	SetScore            bool                    `json:"set_score"`
	Score               int                     `json:"score"`
	SetControversiality bool                    `json:"set_controversiality"`
	Controversiality    float32                 `json:"controversiality"`
	InternalData        MetricQueryInternalData `json:"internal_data"`
}

type RedditSubredditMetricPayload struct {
	ID                 string                  `json:"id"`
	SetSubscribers     bool                    `json:"set_subscribers"`
	Subscribers        int                     `json:"subscribers"`
	SetActiveUserCount bool                    `json:"set_active_user_count"`
	ActiveUserCount    int                     `json:"active_user_count"`
	InternalData       MetricQueryInternalData `json:"internal_data"`
}
