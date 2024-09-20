package temporal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"github.com/jmespath/go-jmespath"
	"go.temporal.io/sdk/log"
	"golang.org/x/sync/errgroup"
	"gonum.org/v1/gonum/stat"
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

func uploadMonitorPosts(l log.Logger, b []byte) error {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("error deserializing response: %w", err)
	}
	// get the number of posts, then iterate over them, extracting each one and
	// creating a corresponding schedule
	iface, err := jmespath.Search("data.children | length(@)", data)
	if err != nil {
		return fmt.Errorf("error extracting post count: %w", err)
	}
	if iface == nil {
		return fmt.Errorf("error extracting post count: nil length")
	}
	count := int(math.Round(iface.(float64)))

	// iterate over the posts and create a schedule for each post. We expect the
	// server to return 409 for most posts, but that's fine it simply means
	// we're already following that post, so just continue on. These run
	// concurrently; as it stands these are taking a sufficiently long time to
	// run serially that it's causing the activity to fail due to start to close
	// timeout errors.
	var errg errgroup.Group
	for i := range count {
		errg.Go(func() error {
			// Skip posts that are stickied but not pinned. We want to track
			// pinned posts, because users will often pin (and sticky) new posts
			// to the top of their profiles. We want to track these. Subreddits
			// will sticky posts as a welcome message, and we don't want to
			// track these.
			iface, err = jmespath.Search(fmt.Sprintf("data.children[%d].data.stickied", i), data)
			if err != nil {
				return fmt.Errorf("error extracting stickied for post %d: %w", i, err)
			}
			if iface == nil {
				return fmt.Errorf("error extracting stickied for post %d: nil stickied", i)
			}
			stickied, ok := iface.(bool)
			if !ok {
				l.Error("sticked assertion error", "stickied", iface)
				stickied = false
			}
			iface, err = jmespath.Search(fmt.Sprintf("data.children[%d].data.pinned", i), data)
			if err != nil {
				return fmt.Errorf("error extracting pinned for post %d: %w", i, err)
			}
			if iface == nil {
				return fmt.Errorf("error extracting pinned for post %d: nil pinned", i)
			}
			pinned, ok := iface.(bool)
			if !ok {
				l.Error("pinned assertion error", "pinned", iface)
				pinned = false
			}
			if stickied && !pinned {
				return nil
			}

			// post id
			iface, err = jmespath.Search(fmt.Sprintf("data.children[%d].data.id", i), data)
			if err != nil {
				return fmt.Errorf("error extracting id for post %d: %w", i, err)
			}
			if iface == nil {
				return fmt.Errorf("error extracting id for post %d: nil id", i)
			}
			id := iface.(string)
			sched := GetDefaultScheduleSpec(RequestKindRedditPost, id)
			payload := api.GenericScheduleRequestPayload{
				RequestKind: RequestKindRedditPost,
				ID:          id,
				Schedule:    sched,
			}
			b, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("error serializing body for post %s: %w", id, err)
			}
			r, err := http.NewRequest(
				http.MethodPost,
				os.Getenv("KAGGO_ENDPOINT")+"/schedule",
				bytes.NewReader(b),
			)
			if err != nil {
				return fmt.Errorf("error creating create schedule request: %w", err)
			}
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				return fmt.Errorf("error doing create schedule request: %w", err)
			}
			defer res.Body.Close()
			// either 200 or 409 means we're good to proceed, short circuit and return early
			if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusConflict {
				return nil
			}
			b, err = io.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("error reading response body for post %s: %w", id, err)
			}
			return fmt.Errorf("bad response (%d) uploading post: %s", res.StatusCode, b)
		})
	}
	return errg.Wait()
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
		Value:        int(math.Round(value)),
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
		Views:        views,
		SetComments:  true,
		Comments:     comments,
		SetLikes:     true,
		Likes:        likes,
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
		Views:          views,
		SetSubscribers: true,
		Subscribers:    subscribers,
		SetVideos:      true,
		Videos:         videos,
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
		Votes:        int(math.Round(votes)),
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
		Views:        int(math.Round(views)),
		SetVotes:     true,
		Votes:        int(math.Round(votes)),
		SetDownloads: true,
		Downloads:    int(math.Round(downloads)),
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
		Score:        int(math.Round(score)),
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
		Score:               int(math.Round(score)),
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
		Subscribers:        int(math.Round(subscribers)),
		SetActiveUserCount: true,
		ActiveUserCount:    int(math.Round(active_user_count)),
		InternalData:       internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetrics(l, "/reddit/subreddit", b)
}

func (a *ActivityRequester) handleRedditSubredditMonitorMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	l.Info(
		"reddit monitor client debug info",
		"remaining", internalData.XRatelimitRemaining,
		"used", internalData.XRatelimitUsed,
		"reset", internalData.XRatelimitReset,
	)
	err := uploadMonitorPosts(l, b)
	if err != nil {
		return nil, fmt.Errorf("error doing subreddit monitor upload: %w", err)
	}
	return &api.DefaultJSONResponse{Message: "ok"}, nil
}

// Handle RequestKindRedditUser requests
func (a *ActivityRequester) handleRedditUserMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}
	// name
	iface, err := jmespath.Search("data.name", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting name: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting name; name is nil")
	}
	name := iface.(string)

	// awardee karma
	iface, err = jmespath.Search("data.awardee_karma", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting awardee_karma: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting awardee_karma; awardee_karma is nil")
	}
	awardee_karma := iface.(float64)

	// awarder karma
	iface, err = jmespath.Search("data.awarder_karma", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting awarder_karma: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting awarder_karma; awarder_karma is nil")
	}
	awarder_karma := iface.(float64)

	// comment karma
	iface, err = jmespath.Search("data.comment_karma", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting comment_karma: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting comment_karma; comment_karma is nil")
	}
	comment_karma := iface.(float64)

	// link karma
	iface, err = jmespath.Search("data.link_karma", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting link_karma: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting link_karma; link_karma is nil")
	}
	link_karma := iface.(float64)

	// total karma
	iface, err = jmespath.Search("data.total_karma", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting total_karma: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting total_karma; name is nil")
	}
	total_karma := iface.(float64)

	// upload the metrics to the server
	payload := api.RedditUserMetricPayload{
		ID:              name, // name is our internal id for reddit users
		SetAwardeeKarma: true,
		AwardeeKarma:    int(awardee_karma),
		SetAwarderKarma: true,
		AwarderKarma:    int(awarder_karma),
		SetCommentKarma: true,
		CommentKarma:    int(comment_karma),
		SetLinkKarma:    true,
		LinkKarma:       int(link_karma),
		SetTotalKarma:   true,
		TotalKarma:      int(total_karma),
		InternalData:    internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetrics(l, "/reddit/user", b)
}

func (a *ActivityRequester) handleRedditUserMonitorMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	err := uploadMonitorPosts(l, b)
	if err != nil {
		return nil, fmt.Errorf("error doing user monitor upload: %w", err)
	}
	return &api.DefaultJSONResponse{Message: "ok"}, nil
}

// Handle RequestKindTwitchClip requests
func (a *ActivityRequester) handleTwitchClipMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
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

	// view_count
	iface, err = jmespath.Search("data[0].view_count", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting view_count: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting view_count; view_count is nil")
	}
	vc := iface.(float64)

	// upload the metrics to the server
	payload := api.TwitchClipMetricPayload{
		ID:           id,
		SetViewCount: true,
		ViewCount:    int(math.Round(vc)),
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetrics(l, "/twitch/clip", b)
}

// Handle RequestKindTwitchVideo requests
func (a *ActivityRequester) handleTwitchVideoMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
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

	// view_count
	iface, err = jmespath.Search("data[0].view_count", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting view_count: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting view_count; view_count is nil")
	}
	vc := iface.(float64)

	// upload the metrics to the server
	payload := api.TwitchVideoMetricPayload{
		ID:           id,
		SetViewCount: true,
		ViewCount:    int(math.Round(vc)),
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetrics(l, "/twitch/video", b)
}

// Handle RequestKindTwitchStream requests
func (a *ActivityRequester) handleTwitchStreamMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}

	// short circuit if nothing is returned; streamer is probably just offline
	iface, err := jmespath.Search("data | length(@)", data)
	if err != nil {
		return nil, fmt.Errorf("error counting stream results: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error counting stream results; count is nil")
	}
	count := iface.(float64)
	if count < 1 {
		return &api.DefaultJSONResponse{Message: "ok"}, nil
	}

	// id
	iface, err = jmespath.Search("data[0].user_login", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting user_login: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting user_login; user_login is nil")
	}
	user_login := iface.(string)

	// view_count
	iface, err = jmespath.Search("data[0].viewer_count", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting viewer_count: %w", err)
	}
	if iface == nil {
		return nil, fmt.Errorf("error extracting viewer_count; viewer_count is nil")
	}
	vc := iface.(float64)

	// upload the metrics to the server
	payload := api.TwitchStreamMetricPayload{
		ID:           user_login,
		SetViewCount: true,
		ViewCount:    int(math.Round(vc)),
		InternalData: internalData,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetrics(l, "/twitch/stream", b)
}

// Handle RequestKindTwitchUserLastDec requests
func (a *ActivityRequester) handleTwitchUserPastDecMetrics(l log.Logger, status int, b []byte, internalData api.MetricQueryInternalData) (*api.DefaultJSONResponse, error) {
	var body struct {
		Data []struct {
			UserID    string `json:"user_login"`
			ViewCount int    `json:"view_count"`
			Duration  string `json:"duration"`
		}
	}

	if err := json.Unmarshal(b, &body); err != nil {
		return nil, fmt.Errorf("error deserializing response: %w", err)
	}

	views := []float64{}
	durs := []float64{}

	for _, v := range body.Data {
		views = append(views, float64(v.ViewCount))
		d, err := time.ParseDuration(v.Duration)
		if err != nil {
			return nil, fmt.Errorf("error parsing video duration %s: %w", v.Duration, err)
		}
		durs = append(durs, float64(d/time.Second))
	}

	slices.Sort(views)
	slices.Sort(durs)

	// views
	avc := stat.Mean(views, nil)
	mvc := stat.Quantile(0.5, stat.LinInterp, views, nil)
	svc := stat.StdDev(views, nil)

	// durations
	ad := stat.Mean(durs, nil)
	md := stat.Quantile(0.5, stat.LinInterp, durs, nil)
	sd := stat.StdDev(durs, nil)

	// upload the metrics to the server
	payload := api.TwitchUserPastDecMetricPayload{
		ID:              body.Data[0].UserID,
		SetAvgViewCount: true,
		AvgViewCount:    float32(avc),
		SetMedViewCount: true,
		MedViewCount:    float32(mvc),
		SetStdViewCount: true,
		StdViewCount:    float32(svc),
		SetAvgDuration:  true,
		AvgDuration:     float32(ad),
		SetMedDuration:  true,
		MedDuration:     float32(md),
		SetStdDuration:  true,
		StdDuration:     float32(sd),
		InternalData:    internalData,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload metadata: %w", err)
	}
	return uploadMetrics(l, "/twitch/user-past-dec", b)
}
