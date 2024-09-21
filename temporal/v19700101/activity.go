package temporal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"github.com/prometheus/client_golang/prometheus"
	"go.temporal.io/sdk/activity"
)

const (
	// request kinds
	RequestKindInternalRandom         = "internal.random"
	RequestKindKaggleNotebook         = "kaggle.notebook"
	RequestKindKaggleDataset          = "kaggle.dataset"
	RequestKindYouTubeVideo           = "youtube.video"
	RequestKindYouTubeChannel         = "youtube.channel"
	RequestKindRedditPost             = "reddit.post"
	RequestKindRedditComment          = "reddit.comment"
	RequestKindRedditSubreddit        = "reddit.subreddit"
	RequestKindRedditSubredditMonitor = "reddit.subreddit-monitor"
	RequestKindRedditUser             = "reddit.user"
	RequestKindRedditUserMonitor      = "reddit.user-monitor"
	RequestKindTwitchClip             = "twitch.clip"
	RequestKindTwitchVideo            = "twitch.video"
	RequestKindTwitchStream           = "twitch.stream"
	RequestKindTwitchUserPastDec      = "twitch.user-past-dec"
	// worker prom metrics
	PromMetricXRatelimitLimit     = "pm-x-ratelimit-limit"
	PromMetricXRatelimitUsed      = "pm-x-ratelimit-used"
	PromMetricXRatelimitRemaining = "pm-x-ratelimit-remaining"
	PromMetricXRatelimitReset     = "pm-x-ratelimit-reset"
)

func GetSupportedRequestKinds() []string {
	return []string{
		RequestKindInternalRandom,
		RequestKindKaggleNotebook,
		RequestKindKaggleDataset,
		RequestKindYouTubeVideo,
		RequestKindYouTubeChannel,
		RequestKindRedditPost,
		RequestKindRedditComment,
		RequestKindRedditSubreddit,
		RequestKindRedditSubredditMonitor,
		RequestKindRedditUser,
		RequestKindRedditUserMonitor,
		RequestKindTwitchClip,
		RequestKindTwitchVideo,
		RequestKindTwitchStream,
		RequestKindTwitchUserPastDec,
	}
}

type ActivityRedditListener struct{}
type ActivityYouTubeListener struct{}

type ActivityRequester struct {
	RedditAuthToken            string
	RedditAuthTokenExp         time.Time
	RedditListenerAuthToken    string
	RedditListenerAuthTokenExp time.Time
	TwitchAuthToken            string
	TwitchAuthTokenExp         time.Time
	Metrics                    map[string]prometheus.Collector
}

func NewActivityRequester() *ActivityRequester {
	pms := map[string]prometheus.Collector{
		PromMetricXRatelimitLimit: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "x_ratelimit_limit",
				Help: "The X-Ratelimit-Limit header from an external server.",
			},
			[]string{"client"},
		),
		PromMetricXRatelimitUsed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "x_ratelimit_used",
				Help: "The X-Ratelimit-Used header from an external server.",
			},
			[]string{"client"},
		),
		PromMetricXRatelimitRemaining: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "x_ratelimit_remaining",
				Help: "The X-Ratelimit-Remaining header from an external server.",
			},
			[]string{"client"},
		),
		PromMetricXRatelimitReset: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "x_ratelimit_reset",
				Help: "The X-Ratelimit-Reset header from an external server.",
			},
			[]string{"client"},
		),
	}
	a := &ActivityRequester{Metrics: pms}
	return a
}

// This is a hook to update requests without updating the originally scheduled
// http.Request. This parses the supplied request and perform any finishing
// touches like setting auth tokens and whatnot. For example, requests to Reddit
// need to have a short lived access token set in the Authorization header. This
// sort frequently changing parameter should not set on the original request
// because the original is hashed to create a unique identifier and prevent
// duplicate schedules.
func (a *ActivityRequester) prepareRequest(drp DoRequestActRequest) (*http.Request, error) {
	buf := bufio.NewReader(bytes.NewReader(drp.Serial))
	r, err := http.ReadRequest(buf)
	if err != nil {
		return nil, fmt.Errorf("error deserializing request: %w", err)
	}

	// https://stackoverflow.com/questions/19595860/http-request-requesturi-field-when-making-request-in-go
	// RequestURI must not be set, but req.URL is incomplete, so parse a new one
	u, err := url.Parse("https://" + r.Host + r.RequestURI)
	if err != nil {
		return nil, fmt.Errorf("error parsing request URL: %s", r.RequestURI)
	}

	r.Header.Set("Accept", "application/json")

	r.URL = u
	r.RequestURI = ""

	switch drp.RequestKind {
	case RequestKindInternalRandom:
		// for internal requests, just set the authorization token
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	case
		RequestKindKaggleNotebook,
		RequestKindKaggleDataset:
		// basic auth
		r.SetBasicAuth(os.Getenv("KAGGLE_USERNAME"), os.Getenv("KAGGLE_API_KEY"))
	case
		RequestKindYouTubeVideo,
		RequestKindYouTubeChannel:
		// for youtube requests, set the non-identifier query params
		q := r.URL.Query()
		q.Set("part", "snippet,contentDetails,statistics")
		q.Set("key", os.Getenv("YOUTUBE_API_KEY"))
		r.URL.RawQuery = q.Encode()

	case
		RequestKindRedditPost,
		RequestKindRedditComment,
		RequestKindRedditSubreddit,
		RequestKindRedditUser:
		// refresh key and set bearer
		err = a.ensureValidRedditToken(time.Duration(60 * time.Second))
		if err != nil {
			return nil, err
		}
		r.Header.Set("User-Agent", os.Getenv("REDDIT_USER_AGENT"))
		r.Header.Set("Authorization", "bearer "+a.RedditAuthToken)
	case
		RequestKindRedditSubredditMonitor,
		RequestKindRedditUserMonitor:
		// refresh key and set bearer
		err = a.ensureValidRedditListenerToken(time.Duration(60 * time.Second))
		if err != nil {
			return nil, err
		}
		r.Header.Set("User-Agent", os.Getenv("REDDIT_LISTENER_USER_AGENT"))
		r.Header.Set("Authorization", "bearer "+a.RedditListenerAuthToken)
		// FIXME: we should set the `after` param here, but we'd need to stick
		// a db cursor onto this struct and start tracking the last monitored
		// post for users and subreddits. Not worth it at the moment.
	case RequestKindTwitchClip, RequestKindTwitchVideo, RequestKindTwitchStream, RequestKindTwitchUserPastDec:
		err = a.ensureValidTwitchToken(time.Duration(60 * time.Second))
		if err != nil {
			return nil, err
		}
		r.Header.Set("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))
		r.Header.Set("Authorization", "Bearer "+a.TwitchAuthToken)
	default:
		return nil, fmt.Errorf("unsupported RequestKind %s", drp.RequestKind)
	}
	return r, nil
}

func (a *ActivityRequester) DoRequest(ctx context.Context, drp DoRequestActRequest) (*DoRequestActResult, error) {
	l := activity.GetLogger(ctx)
	r, err := a.prepareRequest(drp)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// set prometheus metrics; ideally avoid sending the response off to some
	// other caller that may tamper with it
	switch drp.RequestKind {
	case
		RequestKindRedditSubreddit,
		RequestKindRedditUser,
		RequestKindRedditPost,
		RequestKindRedditComment:
		// set X-Ratelimit-Foo headers
		labels := prometheus.Labels{"client": "reddit-poller"}
		a.setRedditPromMetrics(l, labels, resp.Header)
	case
		RequestKindRedditSubredditMonitor,
		RequestKindRedditUserMonitor:
		// set X-Ratelimit-Foo headers
		labels := prometheus.Labels{"client": "reddit-monitor"}
		a.setRedditPromMetrics(l, labels, resp.Header)
	case
		RequestKindTwitchClip,
		RequestKindTwitchVideo,
		RequestKindTwitchStream,
		RequestKindTwitchUserPastDec:
		// set Ratelimit-Foo headers
		labels := prometheus.Labels{"client": "twitch"}
		a.setTwitchPromMetrics(l, labels, resp.Header)
	}

	// return the activity result
	res := DoRequestActResult{
		RequestKind: drp.RequestKind,
		StatusCode:  resp.StatusCode,
		Body:        b,
	}
	return &res, nil
}

// UploadMetadata will handle the response from a get metrics
func (a *ActivityRequester) UploadMetadata(ctx context.Context, drr UploadMetadataActRequest) (*api.DefaultJSONResponse, error) {
	l := activity.GetLogger(ctx)
	switch drr.RequestKind {
	case RequestKindInternalRandom:
		return a.handleInternalRandomMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindYouTubeVideo:
		return a.handleYouTubeVideoMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindYouTubeChannel:
		return a.handleYouTubeChannelMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindKaggleNotebook:
		return a.handleKaggleNotebookMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindKaggleDataset:
		return a.handleKaggleDatasetMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindRedditPost:
		return a.handleRedditPostMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindRedditComment:
		return a.handleRedditCommentMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindRedditSubreddit:
		return a.handleRedditSubredditMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindRedditSubredditMonitor:
		return a.handleRedditSubredditMonitorMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindRedditUser:
		return a.handleRedditUserMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindRedditUserMonitor:
		return a.handleRedditUserMonitorMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindTwitchClip:
		return a.handleTwitchClipMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindTwitchVideo:
		return a.handleTwitchVideoMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindTwitchStream:
		return a.handleTwitchStreamMetadata(l, drr.StatusCode, drr.Body)
	case RequestKindTwitchUserPastDec:
		return a.handleTwitchUserPastDecMetadata(l, drr.StatusCode, drr.Body)
	default:
		return nil, fmt.Errorf("unrecognized RequestKind: %s", drr.RequestKind)
	}
}

// UploadMetrics will handle the response from a get metrics request
func (a *ActivityRequester) UploadMetrics(ctx context.Context, drr UploadMetricsActRequest) (*api.DefaultJSONResponse, error) {
	l := activity.GetLogger(ctx)
	switch drr.RequestKind {
	case RequestKindInternalRandom:
		return a.handleInternalRandomMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindYouTubeVideo:
		return a.handleYouTubeVideoMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindYouTubeChannel:
		return a.handleYouTubeChannelMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindKaggleNotebook:
		return a.handleKaggleNotebookMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindKaggleDataset:
		return a.handleKaggleDatasetMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindRedditPost:
		return a.handleRedditPostMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindRedditComment:
		return a.handleRedditCommentMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindRedditSubreddit:
		return a.handleRedditSubredditMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindRedditSubredditMonitor:
		return a.handleRedditSubredditMonitorMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindRedditUser:
		return a.handleRedditUserMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindRedditUserMonitor:
		return a.handleRedditUserMonitorMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindTwitchClip:
		return a.handleTwitchClipMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindTwitchVideo:
		return a.handleTwitchVideoMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindTwitchStream:
		return a.handleTwitchStreamMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindTwitchUserPastDec:
		return a.handleTwitchUserPastDecMetrics(l, drr.StatusCode, drr.Body)
	default:
		return nil, fmt.Errorf("unrecognized RequestKind: %s", drr.RequestKind)
	}
}
