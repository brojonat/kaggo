package temporal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"go.temporal.io/sdk/activity"
)

const (
	RequestKindInternalRandom = "internal.random"
	RequestKindYouTubeVideo   = "youtube.video"
	RequestKindKaggleNotebook = "kaggle.notebook"
	RequestKindKaggleDataset  = "kaggle.dataset"
	RequestKindRedditPost     = "reddit.post"
	RequestKindRedditComment  = "reddit.comment"
)

type DoRequestParam struct {
	RequestKind string `json:"request_kind"`
	Serial      []byte `json:"serial"`
}
type DoRequestResult struct {
	RequestKind string `json:"request_kind"`
	StatusCode  int    `json:"status_code"`
	Body        []byte `json:"body"`
}

type UploadResponseParam struct {
	ResponseKind string `json:"response_kind"`
	Serial       []byte `json:"serial"`
}

type ActivityRequester struct {
	RedditAuthToken    string
	RedditAuthTokenExp time.Time
}

// Parses the supplied request and perform any finishing touches (for example,
// requests to Reddit need to have a short lived access token set in the
// Authorization header).
func (a *ActivityRequester) prepareRequest(drp DoRequestParam) (*http.Request, error) {
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

	// reddit requests get the auth token
	if strings.HasPrefix(drp.RequestKind, "reddit") {
		// this will refresh the reddit auth token if the deadline is near,
		// otherwise it will just no-op
		a.ensureValidRedditToken(time.Duration(60 * time.Second))
		r.Header.Add("Authorization", "Bearer "+a.RedditAuthToken)
	}

	return r, nil
}

func (a *ActivityRequester) DoRequest(ctx context.Context, drp DoRequestParam) (*DoRequestResult, error) {
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
	res := DoRequestResult{
		RequestKind: drp.RequestKind,
		StatusCode:  resp.StatusCode,
		Body:        b,
	}
	return &res, nil
}

// UploadMetrics will handle the response from a get metrics
func (a *ActivityRequester) UploadMetrics(ctx context.Context, drr DoRequestResult) (*api.DefaultJSONResponse, error) {
	l := activity.GetLogger(ctx)
	switch drr.RequestKind {
	case RequestKindInternalRandom:
		return a.handleInternalRandomMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindYouTubeVideo:
		return a.handleYouTubeVideoMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindKaggleNotebook:
		return a.handleKaggleNotebookMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindKaggleDataset:
		return a.handleKaggleDatasetMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindRedditPost:
		return a.handleRedditPostMetrics(l, drr.StatusCode, drr.Body)
	case RequestKindRedditComment:
		return a.handleRedditCommentMetrics(l, drr.StatusCode, drr.Body)
	default:
		return nil, fmt.Errorf("unrecognized RequestKind: %s", drr.RequestKind)
	}
}
