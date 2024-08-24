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
	"go.temporal.io/sdk/activity"
)

const (
	RequestKindInternalRandom  = "internal.random"
	RequestKindYouTubeVideo    = "youtube.video"
	RequestKindYouTubeChannel  = "youtube.channel"
	RequestKindKaggleNotebook  = "kaggle.notebook"
	RequestKindKaggleDataset   = "kaggle.dataset"
	RequestKindRedditPost      = "reddit.post"
	RequestKindRedditComment   = "reddit.comment"
	RequestKindRedditSubreddit = "reddit.subreddit"
)

type ActivityRequester struct {
	RedditAuthToken    string
	RedditAuthTokenExp time.Time
}

// This is a hook to update requests without updating the originally scheduled
// http.Request. This parses the supplied request and perform any finishing
// touches like setting auth tokens and whatnot. For example, requests to Reddit
// need to have a short lived access token set in the Authorization header.
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
	r.URL = u
	r.RequestURI = ""

	switch drp.RequestKind {
	case RequestKindInternalRandom:
		// nothing to do
	case RequestKindYouTubeVideo, RequestKindYouTubeChannel:
		// re-set the api key
		q := r.URL.Query()
		q.Del("key")
		q.Set("key", os.Getenv("YOUTUBE_API_KEY"))
		r.URL.RawQuery = q.Encode()
	case RequestKindKaggleNotebook:
		// nothing to do
	case RequestKindKaggleDataset:
		// nothing to do
	case RequestKindRedditPost, RequestKindRedditComment, RequestKindRedditSubreddit:
		a.ensureValidRedditToken(time.Duration(60 * time.Second))
		r.Header.Add("Authorization", "Bearer "+a.RedditAuthToken)
	default:
		return nil, fmt.Errorf("unsupported RequestKind %s", drp.RequestKind)
	}
	return r, nil
}

func (a *ActivityRequester) DoRequest(ctx context.Context, drp DoRequestActRequest) (*DoRequestActResult, error) {
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
	default:
		return nil, fmt.Errorf("unrecognized RequestKind: %s", drr.RequestKind)
	}
}
