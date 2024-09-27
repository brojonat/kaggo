package api

import (
	"github.com/brojonat/kaggo/server/db/jsonb"
	"go.temporal.io/sdk/client"
)

type DefaultJSONResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type CreateUserPayload struct {
	Email string                 `json:"email"`
	Data  jsonb.UserMetadataJSON `json:"data"`
}

const (
	UserMetricOpKindAdd         = "add"
	UserMetricOpKindRemove      = "remove"
	UserMetricOpKindAddGroup    = "add-group"
	UserMetricOpKindRemoveGroup = "remove-group"
)

type UserMetricOperationPayload struct {
	OpKind      string `json:"op_kind"`
	Email       string `json:"email"`
	RequestKind string `json:"request_kind"`
	ID          string `json:"id"`
}

type GenericScheduleRequestPayload struct {
	RequestKind string              `json:"request_kind"`
	ID          string              `json:"id"`
	Schedule    client.ScheduleSpec `json:"schedule_spec,omitempty"`
}

type AddListenerSubPayload struct {
	RequestKind string `json:"request_kind"`
	ID          string `json:"id"`
}

type MetricMetadataPayload struct {
	ID          string             `json:"id"`
	RequestKind string             `json:"request_kind"`
	Data        jsonb.MetadataJSON `json:"data"`
}

type InternalMetricPayload struct {
	ID    string `json:"id"`
	Value int    `json:"value"`
}

type YouTubeVideoMetricPayload struct {
	ID          string `json:"id"`
	SetViews    bool   `json:"set_views"`
	Views       int    `json:"views"`
	SetComments bool   `json:"set_comments"`
	Comments    int    `json:"comments"`
	SetLikes    bool   `json:"set_likes"`
	Likes       int    `json:"likes"`
}

type YouTubeChannelMetricPayload struct {
	ID             string `json:"id"`
	SetViews       bool   `json:"set_views"`
	Views          int    `json:"views"`
	SetSubscribers bool   `json:"set_subscribers"`
	Subscribers    int    `json:"subscribers"`
	SetVideos      bool   `json:"set_videos"`
	Videos         int    `json:"videos"`
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

type RedditSubredditMetricPayload struct {
	ID                 string `json:"id"`
	SetSubscribers     bool   `json:"set_subscribers"`
	Subscribers        int    `json:"subscribers"`
	SetActiveUserCount bool   `json:"set_active_user_count"`
	ActiveUserCount    int    `json:"active_user_count"`
}

type RedditUserMetricPayload struct {
	ID              string `json:"id"`
	SetAwardeeKarma bool   `json:"set_awardee_karma"`
	AwardeeKarma    int    `json:"awardee_karma"`
	SetAwarderKarma bool   `json:"set_awarder_karma"`
	AwarderKarma    int    `json:"awarder_karma"`
	SetCommentKarma bool   `json:"set_comment_karma"`
	CommentKarma    int    `json:"comment_karma"`
	SetLinkKarma    bool   `json:"set_like_karma"`
	LinkKarma       int    `json:"like_karma"`
	SetTotalKarma   bool   `json:"set_total_karma"`
	TotalKarma      int    `json:"total_karma"`
}

type TwitchClipMetricPayload struct {
	ID           string `json:"id"`
	SetViewCount bool   `json:"set_view_count"`
	ViewCount    int    `json:"view_count"`
}

type TwitchVideoMetricPayload struct {
	ID           string `json:"id"`
	SetViewCount bool   `json:"set_view_count"`
	ViewCount    int    `json:"view_count"`
}

type TwitchStreamMetricPayload struct {
	ID           string `json:"id"`
	SetViewCount bool   `json:"set_view_count"`
	ViewCount    int    `json:"view_count"`
}

type TwitchUserPastDecMetricPayload struct {
	ID              string  `json:"id"`
	SetAvgViewCount bool    `json:"set_avg_view_count"`
	AvgViewCount    float32 `json:"avg_view_count"`
	SetMedViewCount bool    `json:"set_med_view_count"`
	MedViewCount    float32 `json:"med_view_count"`
	SetStdViewCount bool    `json:"set_std_view_count"`
	StdViewCount    float32 `json:"std_view_count"`
	SetAvgDuration  bool    `json:"set_avg_duration"`
	AvgDuration     float32 `json:"avg_duration"`
	SetMedDuration  bool    `json:"set_med_duration"`
	MedDuration     float32 `json:"med_duration"`
	SetStdDuration  bool    `json:"set_std_duration"`
	StdDuration     float32 `json:"std_duration"`
}
