package temporal

import "github.com/brojonat/kaggo/server/api"

// workflows

type RunYouTubeListenerWFRequest struct{}
type RunRedditListenerWFRequest struct{}

type DoMetadataRequestWFRequest struct {
	RequestKind string `json:"request_kind"`
	Serial      []byte `json:"serial"`
}

type DoPollingRequestWFRequest struct {
	RequestKind string `json:"request_kind"`
	Serial      []byte `json:"serial"`
}

// activities

type YouTubeChannelSubActRequest struct {
	ChannelIDs []string `json:"channel_ids"`
}

type RedditSubActRequest struct {
	Subreddits []string `json:"subreddits"`
	Users      []string `json:"users"`
}

type DoRequestActRequest struct {
	RequestKind string `json:"request_kind"`
	Serial      []byte `json:"serial"`
}
type DoRequestActResult struct {
	RequestKind  string                      `json:"request_kind"`
	StatusCode   int                         `json:"status_code"`
	Body         []byte                      `json:"body"`
	InternalData api.MetricQueryInternalData `json:"internal_data"`
}

type UploadMetadataActRequest struct {
	RequestKind  string                      `json:"request_kind"`
	StatusCode   int                         `json:"status_code"`
	Body         []byte                      `json:"body"`
	InternalData api.MetricQueryInternalData `json:"internal_data"`
}

type UploadMetricsActRequest struct {
	RequestKind  string                      `json:"request_kind"`
	StatusCode   int                         `json:"status_code"`
	Body         []byte                      `json:"body"`
	InternalData api.MetricQueryInternalData `json:"internal_data"`
}
