package temporal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/brojonat/kaggo/server/api"
	"github.com/brojonat/kaggo/server/db/jsonb"
	"github.com/jmespath/go-jmespath"
	"go.temporal.io/sdk/log"
)

// Helper to upload metadata to the kaggo backend
func uploadMetadata(l log.Logger, b []byte) (*api.DefaultJSONResponse, error) {
	endpoint := os.Getenv("KAGGO_ENDPOINT") + "/metadata"
	r, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("error making request to upload metadata: %w", err)
	}
	r.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error doing request to upload metadata: %w", err)
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading metadata upload response body: %w", err)
	}
	var body api.DefaultJSONResponse
	err = json.Unmarshal(b, &body)
	if err != nil {
		l.Error("error parsing metadata upload response body as json", "body", string(b))
		return nil, fmt.Errorf("error parsing metadata upload response: %w", err)
	}
	// do this here after parsing the body so we can access the error message
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("bad response code uploading metadata: %d: %s", res.StatusCode, body.Error)
	}
	return &body, nil
}

// Handle RequestKindInternalRandom requests
func (a *ActivityRequester) handleInternalRandomMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
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

	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindInternalRandom,
		Data:        jsonb.MetadataJSON{ID: id, Link: "https://api.kaggo.brojonat.com/internal/metrics?id=" + id},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindYouTubeVideo requests
func (a *ActivityRequester) handleYouTubeVideoMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
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

	// title
	iface, err = jmespath.Search("items[0].snippet.title", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting title: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting title; title is nil")
	}
	title := iface.(string)

	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindYouTubeVideo,
		Data: jsonb.MetadataJSON{
			ID:    id,
			Link:  fmt.Sprintf("https://www.youtube.com/watch?v=%s", id),
			Title: title,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindYouTubeChannel requests
func (a *ActivityRequester) handleYouTubeChannelMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
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

	// title
	iface, err = jmespath.Search("items[0].snippet.title", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting title: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting title; title is nil")
	}
	title := iface.(string)

	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindYouTubeChannel,
		Data: jsonb.MetadataJSON{
			ID:    id,
			Link:  fmt.Sprintf("https://www.youtube.com/channel/%s", id),
			Title: title,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindKaggleNotebook requests
func (a *ActivityRequester) handleKaggleNotebookMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing notebook response: %w", err)
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
	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindKaggleNotebook,
		Data: jsonb.MetadataJSON{
			ID:   id,
			Link: fmt.Sprintf("https://www.kaggle.com/code/%s", id),
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindKaggleDataset requests
func (a *ActivityRequester) handleKaggleDatasetMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing dataset response: %w", err)
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
	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindKaggleDataset,
		Data: jsonb.MetadataJSON{
			ID:   id,
			Link: fmt.Sprintf("https://www.kaggle.com/datasets/%s", id),
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditPost requests
func (a *ActivityRequester) handleRedditPostMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing post response: %w", err)
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

	// title
	iface, err = jmespath.Search("data.children[0].data.title", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting title: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting title; title is nil")
	}
	title := iface.(string)

	// link
	iface, err = jmespath.Search("data.children[0].data.permalink", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting permalink: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting permalink; id is nil")
	}
	permalink := iface.(string)

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindRedditPost,
		Data: jsonb.MetadataJSON{
			ID:    id,
			Link:  "https://www.reddit.com" + permalink,
			Title: title,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditComment requests
func (a *ActivityRequester) handleRedditCommentMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing comment response: %w", err)
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

	// text
	iface, err = jmespath.Search("data.children[0].data.body", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting body: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting body; body is nil")
	}
	text := iface.(string)

	// author
	iface, err = jmespath.Search("data.children[0].data.author", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting author: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting author; author is nil")
	}
	author := iface.(string)

	// link
	iface, err = jmespath.Search("data.children[0].data.permalink", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting permalink: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting permalink; permalink is nil")
	}
	permalink := iface.(string)

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindRedditComment,
		Data: jsonb.MetadataJSON{
			ID:      id,
			Author:  author,
			Comment: text,
			Link:    "https://www.reddit.com" + permalink,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditSubreddit requests
func (a *ActivityRequester) handleRedditSubredditMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
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

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindRedditSubreddit,
		Data: jsonb.MetadataJSON{
			ID:   id,
			Link: "https://www.reddit.com/r/" + id,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}
