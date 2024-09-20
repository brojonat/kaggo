package temporal

import (
	"time"

	"github.com/brojonat/kaggo/server/api"
	"go.temporal.io/sdk/client"
)

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

func GetDefaultScheduleSpec(rk, id string) client.ScheduleSpec {
	var s client.ScheduleSpec

	// first switch over request kinds to get the base schedule
	switch rk {
	case RequestKindInternalRandom:
		// internal queries are frequent since they're cheap
		s = client.ScheduleSpec{
			Calendars: []client.ScheduleCalendarSpec{
				{
					Second:  []client.ScheduleRange{{Start: 0, End: 59, Step: 30}},
					Minute:  []client.ScheduleRange{{Start: 0, End: 59, Step: 1}},
					Hour:    []client.ScheduleRange{{Start: 0, End: 23, Step: 1}},
					Comment: "every 30 seconds, no jitter",
				},
			},
		}
	case RequestKindYouTubeChannel, RequestKindYouTubeVideo:
		// do youtube queries every 10 minutes; high res isn't super necessary,
		// we have a lot of IDs to query, and the rate limit is pretty much fixed
		s = client.ScheduleSpec{
			Calendars: []client.ScheduleCalendarSpec{
				{
					Second:  []client.ScheduleRange{{Start: 0}},
					Minute:  []client.ScheduleRange{{Start: 0}},
					Hour:    []client.ScheduleRange{{Start: 0, End: 23, Step: 1}},
					Comment: "every hour, with an hour of jitter",
				},
			},
			Jitter: 60 * 60 * 1e9,
		}

	case RequestKindRedditSubredditMonitor, RequestKindRedditUserMonitor:
		// do reddit monitor queries every minute; we want to find posts ASAP, this
		// runs under a different reddit client id and we don't have a ton of ids to
		// monitor
		s = client.ScheduleSpec{
			Calendars: []client.ScheduleCalendarSpec{
				{
					Second:  []client.ScheduleRange{{Start: 0}},
					Minute:  []client.ScheduleRange{{Start: 0, End: 59, Step: 1}},
					Hour:    []client.ScheduleRange{{Start: 0, End: 23, Step: 1}},
					Comment: "every minute, with a minute of jitter",
				},
			},
			Jitter: 60 * 1e9,
		}
	case RequestKindTwitchStream:
		// do twitch stream queries every minute
		s = client.ScheduleSpec{
			Calendars: []client.ScheduleCalendarSpec{
				{
					Second:  []client.ScheduleRange{{Start: 0}},
					Minute:  []client.ScheduleRange{{Start: 0, End: 59, Step: 1}},
					Hour:    []client.ScheduleRange{{Start: 0, End: 23, Step: 1}},
					Comment: "every minute, with a minute of jitter",
				},
			},
			Jitter: 60 * 1e9,
		}
	default:
		// default to every 15 minutes
		s = client.ScheduleSpec{
			Calendars: []client.ScheduleCalendarSpec{
				{
					Second:  []client.ScheduleRange{{Start: 0}},
					Minute:  []client.ScheduleRange{{Start: 0, End: 59, Step: 15}},
					Hour:    []client.ScheduleRange{{Start: 0, End: 23}},
					Comment: "every 15 minutes with 15 minutes of jitter",
				},
			},
			Jitter: 15 * 60 * 1e9,
		}
	}

	// now apply the EndAt depending on the RequestKind
	switch rk {
	case
		// these schedules should run indefinitely; this is the default behavior
		RequestKindInternalRandom,
		RequestKindKaggleNotebook,
		RequestKindKaggleDataset,
		RequestKindRedditSubreddit,
		RequestKindRedditSubredditMonitor,
		RequestKindRedditUser,
		RequestKindRedditUserMonitor,
		RequestKindYouTubeChannel,
		RequestKindTwitchStream,
		RequestKindTwitchUserPastDec:
		// this is a no-op

	case
		// these schedules should run for an intermediate amount of time
		RequestKindRedditPost,
		RequestKindYouTubeVideo,
		RequestKindTwitchVideo:
		// run for 4 weeks
		s.EndAt = time.Now().Add(4 * 7 * 24 * time.Hour)

	case
		// these schedules are relatively short lived and should terminate after
		// a relatively short time
		RequestKindRedditComment,
		RequestKindTwitchClip:
		// run for 1 week
		s.EndAt = time.Now().Add(7 * 24 * time.Hour)

	default:
		// this is a no-op
	}

	return s
}
