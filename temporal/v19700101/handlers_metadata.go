package temporal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
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

// Handle RequestKindInternalRandom metadata requests
func (a *ActivityRequester) handleInternalRandomMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
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
		Data: jsonb.MetadataJSON{
			ID:         id,
			HumanLabel: id,
			Link:       "https://api.kaggo.brojonat.com/internal/metrics?id=" + id,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindKaggleNotebook metadata requests
func (a *ActivityRequester) handleKaggleNotebookMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
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
			ID:         id,
			HumanLabel: id,
			Link:       fmt.Sprintf("https://www.kaggle.com/code/%s", id),
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindKaggleDataset metadata requests
func (a *ActivityRequester) handleKaggleDatasetMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {

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
			ID:         id,
			HumanLabel: id,
			Link:       fmt.Sprintf("https://www.kaggle.com/datasets/%s", id),
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindYouTubeVideo metadata requests
func (a *ActivityRequester) handleYouTubeVideoMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
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

	// channel
	iface, err = jmespath.Search("items[0].snippet.channelTitle", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting channelTitle: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting channelTitle; channelTitle is nil")
	}
	channelTitle := iface.(string)

	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindYouTubeVideo,
		Data: jsonb.MetadataJSON{
			ID:         id,
			HumanLabel: title,
			Link:       fmt.Sprintf("https://www.youtube.com/watch?v=%s", id),
			Title:      title,
			Owner:      channelTitle,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindYouTubeChannel metadata requests
func (a *ActivityRequester) handleYouTubeChannelMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
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
			ID:         id,
			HumanLabel: title,
			Link:       fmt.Sprintf("https://www.youtube.com/channel/%s", id),
			Title:      title,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditPost metadata requests
func (a *ActivityRequester) handleRedditPostMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {

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

	// owner
	iface, err = jmespath.Search("data.children[0].data.author", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting author: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting author; author is nil")
	}
	author := iface.(string)

	// nsfw
	iface, err = jmespath.Search("data.children[0].data.over_18", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting over_18: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting over_18; over_18 is nil")
	}
	nsfw := iface.(bool)

	// tags
	tags := []string{}

	if nsfw {
		tags = append(tags, "NSFW")
	}

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindRedditPost,
		Data: jsonb.MetadataJSON{
			ID:         id,
			HumanLabel: title,
			Link:       "https://www.reddit.com" + permalink,
			Title:      title,
			Owner:      author,
			Tags:       tags,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditComment metadata requests
func (a *ActivityRequester) handleRedditCommentMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {

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

	// link
	iface, err = jmespath.Search("data.children[0].data.permalink", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting permalink: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting permalink; permalink is nil")
	}
	permalink := iface.(string)

	// author
	iface, err = jmespath.Search("data.children[0].data.author", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting author: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting author; author is nil")
	}
	author := iface.(string)

	// text
	iface, err = jmespath.Search("data.children[0].data.body", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting body: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting body; body is nil")
	}
	text := iface.(string)

	// subreddit
	iface, err = jmespath.Search("data.children[0].data.subreddit", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting subreddit: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting subreddit; subreddit is nil")
	}
	subreddit := iface.(string)

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindRedditComment,
		Data: jsonb.MetadataJSON{
			ID:         id,
			HumanLabel: fmt.Sprintf("Comment %s by /u/%s", id, author),
			Owner:      author,
			Comment:    text,
			Link:       "https://www.reddit.com" + permalink,
			Subreddit:  subreddit,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditSubreddit metadata requests
func (a *ActivityRequester) handleRedditSubredditMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
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

	// nsfw
	iface, err = jmespath.Search("data.over18", data) // not a typo, astounding
	if err != nil {
		return nil, fmt.Errorf("error extracting over_18: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting over_18; over_18 is nil")
	}
	nsfw := iface.(bool)

	// tags
	tags := []string{}

	if nsfw {
		tags = append(tags, "NSFW")
	}

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindRedditSubreddit,
		Data: jsonb.MetadataJSON{
			ID:         id,
			HumanLabel: id,
			Link:       "https://www.reddit.com/r/" + id,
			Tags:       tags,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditSubredditMonitor metadata requests
// NOTE: this is a no-op since these requests don't have any pertinent metadata
func (a *ActivityRequester) handleRedditSubredditMonitorMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
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

	// nsfw
	iface, err = jmespath.Search("data.over18", data) // not a typo, astounding
	if err != nil {
		return nil, fmt.Errorf("error extracting over_18: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting over_18; over_18 is nil")
	}
	nsfw := iface.(bool)

	// tags
	tags := []string{}

	if nsfw {
		tags = append(tags, "NSFW")
	}

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindRedditSubredditMonitor,
		Data: jsonb.MetadataJSON{
			ID:         id,
			HumanLabel: id,
			Link:       "https://www.reddit.com/r/" + id,
			Tags:       tags,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditUser metadata requests
func (a *ActivityRequester) handleRedditUserMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}
	// id
	iface, err := jmespath.Search("data.name", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting name: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting name; name is nil")
	}
	name := iface.(string)

	// user_id
	iface, err = jmespath.Search("data.id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting id; id is nil")
	}
	id := iface.(string)

	// created
	iface, err = jmespath.Search("data.created", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting created: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting created; created is nil")
	}
	ts_created := iface.(float64)

	// desc
	iface, err = jmespath.Search("data.subreddit.public_description", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting description: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting description; description is nil")
	}
	desc := iface.(string)

	// nsfw
	iface, err = jmespath.Search("data.subreddit.over_18", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting over_18: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting over_18; over_18 is nil")
	}
	nsfw := iface.(bool)

	// tags
	tags := []string{}

	if nsfw {
		tags = append(tags, "NSFW")
	}

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          name, // name is our internal id for users
		RequestKind: RequestKindRedditUser,
		Data: jsonb.MetadataJSON{
			ID:          name,
			HumanLabel:  name,
			Link:        "https://www.reddit.com/user/" + name,
			Description: desc,
			TSCreated:   int(ts_created),
			UserID:      fmt.Sprintf("t2_%s", id),
			Tags:        tags,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditUserMonitor metadata requests
// NOTE: this is a no-op since these requests don't have any pertinent metadata
func (a *ActivityRequester) handleRedditUserMonitorMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}
	// id
	iface, err := jmespath.Search("data.name", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting name: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting name; name is nil")
	}
	name := iface.(string)

	// user_id
	iface, err = jmespath.Search("data.id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting id; id is nil")
	}
	id := iface.(string)

	// created
	iface, err = jmespath.Search("data.created", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting created: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting created; created is nil")
	}
	ts_created := iface.(float64)

	// desc
	iface, err = jmespath.Search("data.subreddit.public_description", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting description: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting description; description is nil")
	}
	desc := iface.(string)

	// nsfw
	iface, err = jmespath.Search("data.subreddit.over_18", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting over_18: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting over_18; over_18 is nil")
	}
	nsfw := iface.(bool)

	// tags
	tags := []string{}

	if nsfw {
		tags = append(tags, "NSFW")
	}

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          name, // name is our internal id for users
		RequestKind: RequestKindRedditUserMonitor,
		Data: jsonb.MetadataJSON{
			ID:          name,
			HumanLabel:  name,
			Link:        "https://www.reddit.com/user/" + name,
			Description: desc,
			TSCreated:   int(ts_created),
			UserID:      fmt.Sprintf("t2_%s", id),
			Tags:        tags,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindTwitchClip metadata requests
func (a *ActivityRequester) handleTwitchClipMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}
	// id
	iface, err := jmespath.Search("data[0].id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting id; id is nil")
	}
	id := iface.(string)

	// url
	iface, err = jmespath.Search("data[0].url", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting url: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting url; url is nil")
	}
	url := iface.(string)

	// broadcaster_name
	iface, err = jmespath.Search("data[0].broadcaster_name", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting broadcaster_name: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting broadcaster_name; broadcaster_name is nil")
	}
	broadcaster_name := iface.(string)

	// creator_name
	iface, err = jmespath.Search("data[0].creator_name", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting creator_name: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting creator_name; creator_name is nil")
	}
	creator_name := iface.(string)

	// game_id
	iface, err = jmespath.Search("data[0].game_id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting game_id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting game_id; game_id is nil")
	}
	game_id := iface.(string)

	// title
	iface, err = jmespath.Search("data[0].title", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting title: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting title; title is nil")
	}
	title := iface.(string)

	// duration
	iface, err = jmespath.Search("data[0].duration", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting duration: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting duration; duration is nil")
	}
	duration := iface.(float64)

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindTwitchClip,
		Data: jsonb.MetadataJSON{
			ID:          id,
			HumanLabel:  title,
			Link:        url,
			Broadcaster: broadcaster_name,
			Owner:       creator_name,
			GameID:      game_id,
			Title:       title,
			Duration:    int(math.Round(duration)),
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindTwitchVideo metadata requests
func (a *ActivityRequester) handleTwitchVideoMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}
	// id
	iface, err := jmespath.Search("data[0].id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting id; id is nil")
	}
	id := iface.(string)

	// user_name
	iface, err = jmespath.Search("data[0].user_name", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting user_name: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting user_name; user_name is nil")
	}
	user_name := iface.(string)

	// title
	iface, err = jmespath.Search("data[0].title", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting title: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting title; title is nil")
	}
	title := iface.(string)

	// url
	iface, err = jmespath.Search("data[0].url", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting url: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting url; url is nil")
	}
	url := iface.(string)

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindTwitchVideo,
		Data: jsonb.MetadataJSON{
			ID:         id,
			HumanLabel: title,
			Owner:      user_name,
			Title:      title,
			Link:       url,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindTwitchStream metadata requests
func (a *ActivityRequester) handleTwitchStreamMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}
	// user_id
	iface, err := jmespath.Search("data[0].id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting user id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting user id; id is nil")
	}
	user_id := iface.(string)

	// user_login
	iface, err = jmespath.Search("data[0].login", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting user login: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting user login; login is nil")
	}
	user_login := iface.(string)

	// user_name
	iface, err = jmespath.Search("data[0].display_name", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting display_name: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting display_name; display_name is nil")
	}
	display_name := iface.(string)

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		// NOTE: the ID field for this type of request is the user slug, NOT the user ID because
		// twitch makes it nearly impossible to find the user_id
		ID:          user_login,
		RequestKind: RequestKindTwitchStream,
		Data: jsonb.MetadataJSON{
			ID:          user_login,
			HumanLabel:  display_name,
			Owner:       user_login,
			Link:        "https://twitch.tv/" + user_login,
			DisplayName: display_name,
			UserID:      user_id,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)

}

// Handle RequestKindTwitchUserPastDec metadata requests
func (a *ActivityRequester) handleTwitchUserPastDecMetadata(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}
	// user_id
	iface, err := jmespath.Search("data[0].id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting user id: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting user id; id is nil")
	}
	user_id := iface.(string)

	// user_login
	iface, err = jmespath.Search("data[0].login", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting user login: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting user login; login is nil")
	}
	user_login := iface.(string)

	// user_name
	iface, err = jmespath.Search("data[0].display_name", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting display_name: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting display_name; display_name is nil")
	}
	display_name := iface.(string)

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		// NOTE: the ID field for this type of request is the user slug, NOT the user ID because
		// twitch makes it nearly impossible to find the user_id
		ID:          user_login,
		RequestKind: RequestKindTwitchUserPastDec,
		Data: jsonb.MetadataJSON{
			ID:          user_login,
			HumanLabel:  display_name,
			Owner:       user_login,
			Link:        "https://twitch.tv/" + user_login,
			DisplayName: display_name,
			UserID:      user_id,
		},
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetadata(l, b)
}
