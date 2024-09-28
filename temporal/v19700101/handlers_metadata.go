package temporal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

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
		Data: jsonb.MetadataJSON{
			ID:         id,
			HumanLabel: id,
			Link:       "https://api.kaggo.brojonat.com/internal/metrics?id=" + id,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindKaggleNotebook metadata requests
func (a *ActivityRequester) handleKaggleNotebookMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing notebook response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("[0].ref", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting ref: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting ref; ref is nil")}
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
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindKaggleDataset metadata requests
func (a *ActivityRequester) handleKaggleDatasetMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing dataset response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("[0].ref", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting ref: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting ref; ref is nil")}
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
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindYouTubeVideo metadata requests
func (a *ActivityRequester) handleYouTubeVideoMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing internal response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("items[0].id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id; id is nil")}
	}
	id := iface.(string)

	// title
	iface, err = jmespath.Search("items[0].snippet.title", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting title: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting title; title is nil")}
	}
	title := iface.(string)

	// timestamp
	iface, err = jmespath.Search("items[0].snippet.publishedAt", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting publishedAt: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting publishedAt; publishedAt is nil")}
	}
	tsRaw := iface.(string)
	ts, err := time.Parse(time.RFC3339, tsRaw)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error parsing publishedAt: %w", err)}
	}

	// channel
	iface, err = jmespath.Search("items[0].snippet.channelId", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting channelId: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting channelId; channelId is nil")}
	}
	channelID := iface.(string)

	iface, err = jmespath.Search("items[0].snippet.channelTitle", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting channelTitle: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting channelTitle; channelTitle is nil")}
	}
	channelTitle := iface.(string)

	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindYouTubeVideo,
		Data: jsonb.MetadataJSON{
			ID:                 id,
			HumanLabel:         title,
			Link:               fmt.Sprintf("https://www.youtube.com/watch?v=%s", id),
			TSCreated:          ts,
			Title:              title,
			ParentChannelID:    channelID,
			ParentChannelTitle: channelTitle,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload data: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindYouTubeChannel metadata requests
func (a *ActivityRequester) handleYouTubeChannelMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing internal response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("items[0].id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id; id is nil")}
	}
	id := iface.(string)

	// title
	iface, err = jmespath.Search("items[0].snippet.title", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting title: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting title; title is nil")}
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
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload data: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditPost metadata requests
func (a *ActivityRequester) handleRedditPostMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing post response: %w", err)}
	}

	// id
	iface, err := jmespath.Search("data.children[0].data.id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id; id is nil")}
	}
	id := iface.(string)

	// title
	iface, err = jmespath.Search("data.children[0].data.title", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting title: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting title; title is nil")}
	}
	title := iface.(string)

	// created FIXME: validate
	iface, err = jmespath.Search("data.children[0].data.created", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting created: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting created; created is nil")}
	}
	ts_unix := iface.(float64)
	ts := time.Unix(int64(math.Round(ts_unix)), 0)

	// link
	iface, err = jmespath.Search("data.children[0].data.permalink", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting permalink: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting permalink; id is nil")}
	}
	permalink := iface.(string)

	// author name
	iface, err = jmespath.Search("data.children[0].data.author", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting author: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting author; author is nil")}
	}
	author_name := iface.(string)

	// author_id may not be available in many cases (e.g., if post deleted or crossposted),
	// so just set it if we can (don't worry about returning an error here)
	author_id := ""
	iface, err = jmespath.Search("data.children[0].data.author_fullname", data)
	if err != nil {
		l.Info(fmt.Sprintf("error extracting author_fullname: %s", err.Error()))
	} else if iface == nil {
		l.Info("error extracting author_fullname; author_fullname is nil")
	} else {
		author_id = iface.(string)
	}

	// subreddit
	iface, err = jmespath.Search("data.children[0].data.subreddit", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting subreddit: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting subreddit; subreddit is nil")}
	}
	subreddit := iface.(string)

	// nsfw
	iface, err = jmespath.Search("data.children[0].data.over_18", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting over_18: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting over_18; over_18 is nil")}
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
			ID:              id,
			HumanLabel:      title,
			Link:            "https://www.reddit.com" + permalink,
			TSCreated:       ts,
			Title:           title,
			ParentUserID:    author_id,
			ParentUserName:  author_name,
			ParentSubreddit: subreddit,
			Tags:            tags,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditComment metadata requests
func (a *ActivityRequester) handleRedditCommentMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing comment response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("data.children[0].data.id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id; id is nil")}
	}
	id := iface.(string)

	// created
	iface, err = jmespath.Search("data.children[0].data.created", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting created: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting created; created is nil")}
	}
	ts_unix := iface.(float64)
	ts := time.Unix(int64(math.Round(ts_unix)), 0)

	// link
	iface, err = jmespath.Search("data.children[0].data.permalink", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting permalink: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting permalink; permalink is nil")}
	}
	permalink := iface.(string)

	// link to parent post
	iface, err = jmespath.Search("data.children[0].data.link_id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting link_id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting link_id; link_id is nil")}
	}
	post_id := iface.(string)

	// name of the parent post from permalink
	// e.g. "/r/{subreddit}/comments/{post-fullname}/{post-title}/{comment-fullname}/",
	post_title := strings.Split(permalink, "/")[4]

	// author
	iface, err = jmespath.Search("data.children[0].data.author", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting author: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting author; author is nil")}
	}
	author_name := iface.(string)

	// author_id
	// author_id may not be available in many cases (e.g., if post deleted or crossposted),
	// so just set it if we can (don't worry about returning an error here)
	author_id := ""
	iface, err = jmespath.Search("data.children[0].data.author_fullname", data)
	if err != nil {
		l.Info(fmt.Sprintf("error extracting author_fullname: %s", err.Error()))
	} else if iface == nil {
		l.Info("error extracting author_fullname; author_fullname is nil")
	} else {
		author_id = iface.(string)
	}

	// text
	iface, err = jmespath.Search("data.children[0].data.body", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting body: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting body; body is nil")}
	}
	text := iface.(string)

	// subreddit
	iface, err = jmespath.Search("data.children[0].data.subreddit", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting subreddit: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting subreddit; subreddit is nil")}
	}
	subreddit := iface.(string)

	// upload the metadata to the server
	payload := api.MetricMetadataPayload{
		ID:          id,
		RequestKind: RequestKindRedditComment,
		Data: jsonb.MetadataJSON{
			ID:              id,
			HumanLabel:      fmt.Sprintf("Comment %s by /u/%s", id, author_name),
			TSCreated:       ts,
			ParentUserID:    author_id,
			ParentUserName:  author_name,
			ParentPostID:    post_id,
			ParentPostTitle: post_title,
			ParentSubreddit: subreddit,
			Comment:         text,
			Link:            "https://www.reddit.com" + permalink,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditSubreddit metadata requests
func (a *ActivityRequester) handleRedditSubredditMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("data.display_name", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id; id is nil")}
	}
	id := iface.(string)

	// nsfw
	iface, err = jmespath.Search("data.over18", data) // not a typo, astounding
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting over_18: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting over_18; over_18 is nil")}
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
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditSubredditMonitor metadata requests
// NOTE: this is a no-op since these requests don't have any pertinent metadata
func (a *ActivityRequester) handleRedditSubredditMonitorMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("data.display_name", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id; id is nil")}
	}
	id := iface.(string)

	// nsfw
	iface, err = jmespath.Search("data.over18", data) // not a typo, astounding
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting over_18: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting over_18; over_18 is nil")}
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
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditUser metadata requests
func (a *ActivityRequester) handleRedditUserMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("data.name", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting name: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting name; name is nil")}
	}
	name := iface.(string)

	// user_id
	iface, err = jmespath.Search("data.id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id; id is nil")}
	}
	id := iface.(string)

	// created
	iface, err = jmespath.Search("data.created", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting created: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting created; created is nil")}
	}
	ts_unix := iface.(float64)
	ts := time.Unix(int64(math.Round(ts_unix)), 0)

	// desc
	iface, err = jmespath.Search("data.subreddit.public_description", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting description: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting description; description is nil")}
	}
	desc := iface.(string)

	// nsfw
	iface, err = jmespath.Search("data.subreddit.over_18", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting over_18: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting over_18; over_18 is nil")}
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
			TSCreated:   ts,
			UserID:      fmt.Sprintf("t2_%s", id),
			Tags:        tags,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindRedditUserMonitor metadata requests
// NOTE: this is a no-op since these requests don't have any pertinent metadata
func (a *ActivityRequester) handleRedditUserMonitorMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("data.name", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting name: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting name; name is nil")}
	}
	name := iface.(string)

	// user_id
	iface, err = jmespath.Search("data.id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id; id is nil")}
	}
	id := iface.(string)

	// created
	iface, err = jmespath.Search("data.created", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting created: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting created; created is nil")}
	}
	ts_unix := iface.(float64)
	ts := time.Unix(int64(math.Round(ts_unix)), 0)

	// desc
	iface, err = jmespath.Search("data.subreddit.public_description", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting description: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting description; description is nil")}
	}
	desc := iface.(string)

	// nsfw
	iface, err = jmespath.Search("data.subreddit.over_18", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting over_18: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting over_18; over_18 is nil")}
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
			TSCreated:   ts,
			UserID:      fmt.Sprintf("t2_%s", id),
			Tags:        tags,
		},
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindTwitchClip metadata requests
func (a *ActivityRequester) handleTwitchClipMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("data[0].id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id; id is nil")}
	}
	id := iface.(string)

	// url
	iface, err = jmespath.Search("data[0].url", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting url: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting url; url is nil")}
	}
	url := iface.(string)

	// broadcaster_name
	iface, err = jmespath.Search("data[0].broadcaster_name", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting broadcaster_name: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting broadcaster_name; broadcaster_name is nil")}
	}
	broadcaster_name := iface.(string)

	// creator_name
	iface, err = jmespath.Search("data[0].creator_name", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting creator_name: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting creator_name; creator_name is nil")}
	}
	creator_name := iface.(string)

	// game_id
	iface, err = jmespath.Search("data[0].game_id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting game_id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting game_id; game_id is nil")}
	}
	game_id := iface.(string)

	// title
	iface, err = jmespath.Search("data[0].title", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting title: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting title; title is nil")}
	}
	title := iface.(string)

	// duration
	iface, err = jmespath.Search("data[0].duration", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting duration: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting duration; duration is nil")}
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
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindTwitchVideo metadata requests
func (a *ActivityRequester) handleTwitchVideoMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing response: %w", err)}
	}
	// id
	iface, err := jmespath.Search("data[0].id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting id; id is nil")}
	}
	id := iface.(string)

	// user_name
	iface, err = jmespath.Search("data[0].user_name", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting user_name: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting user_name; user_name is nil")}
	}
	user_name := iface.(string)

	// title
	iface, err = jmespath.Search("data[0].title", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting title: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting title; title is nil")}
	}
	title := iface.(string)

	// url
	iface, err = jmespath.Search("data[0].url", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting url: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting url; url is nil")}
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
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}

// Handle RequestKindTwitchStream metadata requests
func (a *ActivityRequester) handleTwitchStreamMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing response: %w", err)}
	}
	// user_id
	iface, err := jmespath.Search("data[0].id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting user id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting user id; id is nil")}
	}
	user_id := iface.(string)

	// user_login
	iface, err = jmespath.Search("data[0].login", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting user login: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting user login; login is nil")}
	}
	user_login := iface.(string)

	// user_name
	iface, err = jmespath.Search("data[0].display_name", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting display_name: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting display_name; display_name is nil")}
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
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)

}

// Handle RequestKindTwitchUserPastDec metadata requests
func (a *ActivityRequester) handleTwitchUserPastDecMetadata(l log.Logger, status int, b []byte) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error deserializing response: %w", err)}
	}
	// user_id
	iface, err := jmespath.Search("data[0].id", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting user id: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting user id; id is nil")}
	}
	user_id := iface.(string)

	// user_login
	iface, err = jmespath.Search("data[0].login", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting user login: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting user login; login is nil")}
	}
	user_login := iface.(string)

	// user_name
	iface, err = jmespath.Search("data[0].display_name", data)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting display_name: %w", err)}
	}
	if iface == nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error extracting display_name; display_name is nil")}
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
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, ErrNoRetry{Err: fmt.Errorf("error serializing upload metadata: %w", err)}
	}
	return uploadMetadata(l, b)
}
