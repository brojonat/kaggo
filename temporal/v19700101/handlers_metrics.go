package temporal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/brojonat/kaggo/server/api"
	"github.com/jmespath/go-jmespath"
	"go.temporal.io/sdk/log"
)

// Helper to upload metrics to the kaggo backend
func uploadMetrics(l log.Logger, path string, b []byte) (*api.DefaultJSONResponse, error) {
	endpoint := os.Getenv("KAGGO_ENDPOINT") + path
	r, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("error making request to upload metrics: %w", err)
	}
	r.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error doing request to upload metrics: %w", err)
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading metrics response body: %w", err)
	}
	var body api.DefaultJSONResponse
	err = json.Unmarshal(b, &body)
	if err != nil {
		l.Error("error parsing metrics response body as json", "body", string(b))
		return nil, fmt.Errorf("error parsing metrics response: %w", err)
	}
	// do this here after parsing the body so we can access the error message
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("bad response code uploading metrics: %d: %s", res.StatusCode, body.Error)
	}
	return &body, nil
}

// Handle RequestKindInternalRandom requests
func (a *ActivityRequester) handleInternalRandomMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing internal response: %w", err)
	}

	// id
	iface, err := jmespath.Search("id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting id; id is nil")
	}
	id := iface.(string)

	// value
	iface, err = jmespath.Search("value", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting value: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting value; value is nil")
	}
	value := iface.(float64)

	payload := api.InternalMetricPayload{
		ID:           id,
		Value:        int(value),
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/internal/metrics", b)
}

// Handle RequestKindYouTubeVideo requests
func (a *ActivityRequester) handleYouTubeVideoMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}

	// id
	iface, err := jmespath.Search("items[0].id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting id; id is nil")
	}
	id := iface.(string)

	// views
	iface, err = jmespath.Search("items[0].statistics.viewCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting views: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting views; views is nil")
	}
	viewsStr := iface.(string)
	views, err := strconv.Atoi(viewsStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing views")
	}

	// comments
	iface, err = jmespath.Search("items[0].statistics.commentCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting commentCount: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting commentCount; commentCount is nil")
	}
	commentsStr := iface.(string)
	comments, err := strconv.Atoi(commentsStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing comments")
	}

	// likes
	iface, err = jmespath.Search("items[0].statistics.likeCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting likeCount: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting likeCount; likeCount is nil")
	}
	likesStr := iface.(string)
	likes, err := strconv.Atoi(likesStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing likes")
	}

	payload := api.YouTubeVideoMetricPayload{
		ID:           id,
		SetViews:     true,
		Views:        int(views),
		SetComments:  true,
		Comments:     int(comments),
		SetLikes:     true,
		Likes:        int(likes),
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/youtube/video", b)
}

// Handle RequestKindYouTubeChannel requests
func (a *ActivityRequester) handleYouTubeChannelMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing internal response: %w", err)
	}

	// id
	iface, err := jmespath.Search("items[0].id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting id; id is nil")
	}
	id := iface.(string)

	// views
	iface, err = jmespath.Search("items[0].statistics.viewCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting views: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting views; views is nil")
	}
	viewsStr := iface.(string)
	views, err := strconv.Atoi(viewsStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing views")
	}

	// subscribers
	iface, err = jmespath.Search("items[0].statistics.subscriberCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting subscriberCount: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting subscriberCount; subscriberCount is nil")
	}
	subStr := iface.(string)
	subscribers, err := strconv.Atoi(subStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing subscribers")
	}

	// videos
	iface, err = jmespath.Search("items[0].statistics.videoCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting videoCount: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting videoCount; videoCount is nil")
	}
	videosStr := iface.(string)
	videos, err := strconv.Atoi(videosStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing videos")
	}

	payload := api.YouTubeChannelMetricPayload{
		ID:             id,
		SetViews:       true,
		Views:          int(views),
		SetSubscribers: true,
		Subscribers:    int(subscribers),
		SetVideos:      true,
		Videos:         int(videos),
		InternalData:   internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/youtube/channel", b)
}

// Handle RequestKindKaggleNotebook requests
func (a *ActivityRequester) handleKaggleNotebookMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}

	// id
	iface, err := jmespath.Search("[0].ref", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting ref: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting ref; ref is nil")
	}
	id := iface.(string)

	// votes
	iface, err = jmespath.Search("[0].totalVotes", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting votes: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting totalVotes; totalVotes is nil")
	}
	votes := iface.(float64)

	// upload the metrics to the server
	payload := api.KaggleNotebookMetricPayload{
		ID:           id,
		SetVotes:     true,
		Votes:        int(votes),
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/kaggle/notebook", b)
}

// Handle RequestKindKaggleDataset requests
func (a *ActivityRequester) handleKaggleDatasetMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}

	// id
	iface, err := jmespath.Search("[0].ref", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting ref: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting ref; ref is nil")
	}
	id := iface.(string)

	// views
	iface, err = jmespath.Search("[0].viewCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting viewCount: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting viewCount; viewCount is nil")
	}
	views := iface.(float64)

	// votes
	iface, err = jmespath.Search("[0].voteCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting votes: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting voteCount; voteCount is nil")
	}
	votes := iface.(float64)

	// downloads
	iface, err = jmespath.Search("[0].downloadCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting downloadCount: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting downloadCount; downloadCount is nil")
	}
	downloads := iface.(float64)

	// upload the metrics to the server
	payload := api.KaggleDatasetMetricPayload{
		ID:           id,
		SetViews:     true,
		Views:        int(views),
		SetVotes:     true,
		Votes:        int(votes),
		SetDownloads: true,
		Downloads:    int(downloads),
		InternalData: internalData,
	}

	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/kaggle/dataset", b)
}

// Handle RequestKindRedditPost requests
func (a *ActivityRequester) handleRedditPostMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}

	// id
	iface, err := jmespath.Search("data.children[0].data.id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting id; id is nil")
	}
	id := iface.(string)

	// score
	iface, err = jmespath.Search("data.children[0].data.score", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting score: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting score; score is nil")
	}
	score := iface.(float64)

	// ratio
	iface, err = jmespath.Search("data.children[0].data.upvote_ratio", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting upvote_ratio: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting upvote_ratio; upvote_ratio is nil")
	}
	ratio := iface.(float64)

	// upload the metrics to the server
	payload := api.RedditPostMetricPayload{
		ID:           id,
		SetScore:     true,
		Score:        int(score),
		SetRatio:     true,
		Ratio:        float32(ratio),
		InternalData: internalData,
	}

	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/reddit/post", b)
}

// Handle RequestKindRedditComment requests
func (a *ActivityRequester) handleRedditCommentMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}

	// id
	iface, err := jmespath.Search("data.children[0].data.id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting id; id is nil")
	}
	id := iface.(string)

	// score
	iface, err = jmespath.Search("data.children[0].data.score", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting score: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting score; score is nil")
	}
	score := iface.(float64)

	// controversiality
	iface, err = jmespath.Search("data.children[0].data.controversiality", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting controversiality; controversiality is nil")
	}
	cont := iface.(float64)

	// upload the metrics to the server
	payload := api.RedditCommentMetricPayload{
		ID:                  id,
		SetScore:            true,
		Score:               int(score),
		SetControversiality: true,
		Controversiality:    float32(cont),
		InternalData:        internalData,
	}

	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/reddit/comment", b)
}

// Handle RequestKindRedditSubreddit requests
func (a *ActivityRequester) handleRedditSubredditMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}
	// id
	iface, err := jmespath.Search("data.display_name", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting id; id is nil")
	}
	id := iface.(string)

	// subscribers
	iface, err = jmespath.Search("data.subscribers", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting subscribers: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting subscribers; subscribers is nil")
	}
	subscribers := iface.(float64)

	// active user count
	iface, err = jmespath.Search("data.active_user_count", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting active_user_count: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting active_user_count; active_user_count is nil")
	}
	active_user_count := iface.(float64)

	// upload the metrics to the server
	payload := api.RedditSubredditMetricPayload{
		ID:                 id,
		SetSubscribers:     true,
		Subscribers:        int(subscribers),
		SetActiveUserCount: true,
		ActiveUserCount:    int(active_user_count),
		InternalData:       internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetrics(l, "/reddit/subreddit", b)
}
