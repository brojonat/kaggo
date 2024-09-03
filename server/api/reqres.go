package api

import (
	"github.com/brojonat/kaggo/server/db/jsonb"
	"github.com/turnage/graw/reddit"
	"go.temporal.io/sdk/client"
)

func GetDefaultScheduleSpec(rk, id string) client.ScheduleSpec {
	s := client.ScheduleSpec{
		Calendars: []client.ScheduleCalendarSpec{
			{
				Second:  []client.ScheduleRange{{Start: 0}},
				Minute:  []client.ScheduleRange{{Start: 0, End: 59, Step: 5}},
				Hour:    []client.ScheduleRange{{Start: 0, End: 23}},
				Comment: "every 5 minutes",
			},
		},
		Jitter: 300000000000,
	}
	return s
}

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

// Contains arbitrary data that metric (or metadata) handlers want to pass back to the server.
type MetricQueryInternalData struct {
	// used by reddit
	XRatelimitUsed      string `json:"x_ratelimit_used"`
	XRatelimitRemaining string `json:"x_ratelimit_remaining"`
	XRatelimitReset     string `json:"x_ratelimit_reset"`
	// used by twitch
	RatelimitLimit     string `json:"ratelimit_limit"`
	RatelimitRemaining string `json:"ratelimit_remaining"`
	RatelimitReset     string `json:"ratelimit_reset"`
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

type TwitchClipMetricPayload struct {
	ID           string                  `json:"id"`
	SetViewCount bool                    `json:"set_view_count"`
	ViewCount    int                     `json:"view_count"`
	InternalData MetricQueryInternalData `json:"internal_data"`
}

type TwitchVideoMetricPayload struct {
	ID           string                  `json:"id"`
	SetViewCount bool                    `json:"set_view_count"`
	ViewCount    int                     `json:"view_count"`
	InternalData MetricQueryInternalData `json:"internal_data"`
}

type TwitchStreamMetricPayload struct {
	ID           string                  `json:"id"`
	SetViewCount bool                    `json:"set_view_count"`
	ViewCount    int                     `json:"view_count"`
	InternalData MetricQueryInternalData `json:"internal_data"`
}

type TwitchUserPastDecMetricPayload struct {
	ID              string                  `json:"id"`
	SetAvgViewCount bool                    `json:"set_avg_view_count"`
	AvgViewCount    float32                 `json:"avg_view_count"`
	SetMedViewCount bool                    `json:"set_med_view_count"`
	MedViewCount    float32                 `json:"med_view_count"`
	SetStdViewCount bool                    `json:"set_std_view_count"`
	StdViewCount    float32                 `json:"std_view_count"`
	SetAvgDuration  bool                    `json:"set_avg_duration"`
	AvgDuration     float32                 `json:"avg_duration"`
	SetMedDuration  bool                    `json:"set_med_duration"`
	MedDuration     float32                 `json:"med_duration"`
	SetStdDuration  bool                    `json:"set_std_duration"`
	StdDuration     float32                 `json:"std_duration"`
	InternalData    MetricQueryInternalData `json:"internal_data"`
}

type RedditPostUpdate struct {
	Post reddit.Post `json:"post"`
}
