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
func (a *ActivityRequester) handleInternalRandomMetrics(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing internal response: %w", err)
	}

	// id
	iface, err := jmespath.Search("id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting internal id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting internal id; id is nil")
	}
	id := iface.(string)

	// value
	iface, err = jmespath.Search("value", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting internal value: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting internal value; value is nil")
	}
	value := iface.(float64)

	payload := api.InternalMetricPayload{
		ID:    id,
		Value: int(value),
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/internal/metrics", b)
}

// Handle RequestKindYouTubeVideo requests
func (a *ActivityRequester) handleYouTubeVideoMetrics(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
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
		return nil, fmt.Errorf("error extracting video id; id is nil")
	}
	id := iface.(string)

	// views
	iface, err = jmespath.Search("items[0].statistics.viewCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting views: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting video views; views is nil")
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
		return nil, fmt.Errorf("error extracting video commentCount; commentCount is nil")
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
		return nil, fmt.Errorf("error extracting video likeCount; likeCount is nil")
	}
	likesStr := iface.(string)
	likes, err := strconv.Atoi(likesStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing likes")
	}

	payload := api.YouTubeVideoMetricPayload{
		ID:          id,
		SetViews:    true,
		Views:       int(views),
		SetComments: true,
		Comments:    int(comments),
		SetLikes:    true,
		Likes:       int(likes),
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/youtube/video", b)
}

// Handle RequestKindKaggleNotebook requests
func (a *ActivityRequester) handleKaggleNotebookMetrics(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing notebook response: %w", err)
	}

	// id
	iface, err := jmespath.Search("[0].ref", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting notebook ref: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting notebook ref; ref is nil")
	}
	id := iface.(string)

	// votes
	iface, err = jmespath.Search("[0].totalVotes", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting notebook votes: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting notebook totalVotes; totalVotes is nil")
	}
	votes := iface.(float64)

	// upload the metrics to the server
	payload := api.KaggleNotebookMetricPayload{
		ID:       id,
		SetVotes: true,
		Votes:    int(votes),
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/kaggle/notebook", b)
}

// Handle RequestKindKaggleDataset requests
func (a *ActivityRequester) handleKaggleDatasetMetrics(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing dataset response: %w", err)
	}

	// id
	iface, err := jmespath.Search("[0].ref", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting dataset ref: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting dataset ref; ref is nil")
	}
	id := iface.(string)

	// views
	iface, err = jmespath.Search("[0].viewCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting dataset viewCount: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting dataset viewCount; viewCount is nil")
	}
	views := iface.(float64)

	// votes
	iface, err = jmespath.Search("[0].voteCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting dataset votes: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting dataset voteCount; voteCount is nil")
	}
	votes := iface.(float64)

	// downloads
	iface, err = jmespath.Search("[0].downloadCount", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting dataset downloadCount: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting dataset downloadCount; downloadCount is nil")
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
	}

	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/kaggle/dataset", b)
}

// Handle RequestKindRedditPost requests
func (a *ActivityRequester) handleRedditPostMetrics(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing dataset response: %w", err)
	}

	// id
	iface, err := jmespath.Search("data.children[0].data.id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting post id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting post id; id is nil")
	}
	id := iface.(string)

	// score
	iface, err = jmespath.Search("data.children[0].data.score", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting post score: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting post score; score is nil")
	}
	score := iface.(float64)

	// ratio
	iface, err = jmespath.Search("data.children[0].data.upvote_ratio", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting post upvote_ratio: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting post upvote_ratio; upvote_ratio is nil")
	}
	ratio := iface.(float64)

	// upload the metrics to the server
	payload := api.RedditPostMetricPayload{
		ID:       id,
		SetScore: true,
		Score:    int(score),
		SetRatio: true,
		Ratio:    float32(ratio),
	}

	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/reddit/post", b)
}

// Handle RequestKindRedditComment requests
func (a *ActivityRequester) handleRedditCommentMetrics(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing dataset response: %w", err)
	}

	// id
	iface, err := jmespath.Search("data.children[0].data.id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting comment id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting comment id; id is nil")
	}
	id := iface.(string)

	// score
	iface, err = jmespath.Search("data.children[0].data.score", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting comment score: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting comment score; score is nil")
	}
	score := iface.(float64)

	// controversiality
	iface, err = jmespath.Search("data.children[0].data.controversiality", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting comment id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting comment controversiality; controversiality is nil")
	}
	cont := iface.(float64)

	// upload the metrics to the server
	payload := api.RedditCommentMetricPayload{
		ID:                  id,
		SetScore:            true,
		Score:               int(score),
		SetControversiality: true,
		Controversiality:    float32(cont),
	}

	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetrics(l, "/reddit/comment", b)
}
