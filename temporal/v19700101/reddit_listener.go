package temporal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/brojonat/kaggo/server/api"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
	"go.temporal.io/sdk/activity"
)

type redditHandler struct{}

func (rh *redditHandler) validatePost(p *reddit.Post) error {
	// skip stickied posts because they always show up in the crawler
	if p.Stickied {
		return fmt.Errorf("skip stickied post")
	}
	return nil
}

func (rh *redditHandler) TearDown() {}
func (rh *redditHandler) Post(p *reddit.Post) error {

	// if the post isn't a good post, just short circuit here
	if rh.validatePost(p) != nil {
		return nil
	}

	body := api.RedditPostUpdate{
		Post: *p,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, os.Getenv("KAGGO_ENDPOINT")+"/notification/reddit/post", bytes.NewReader(b))
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %s", res.Status)
	}
	return nil
}
func (rh *redditHandler) Comment(post *reddit.Comment) error    { return nil }
func (rh *redditHandler) Message(msg *reddit.Message) error     { return nil }
func (rh *redditHandler) PostReply(reply *reddit.Message) error { return nil }
func (rh *redditHandler) UserPost(p *reddit.Post) error {

	// if the post isn't a good post, just short circuit here
	if rh.validatePost(p) != nil {
		return nil
	}

	body := api.RedditPostUpdate{
		Post: *p,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, os.Getenv("KAGGO_ENDPOINT")+"/notification/reddit/post", bytes.NewReader(b))
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %s", res.Status)
	}
	return nil
}
func (rh *redditHandler) UserComment(c *reddit.Comment) error      { return nil }
func (rh *redditHandler) CommentReply(reply *reddit.Message) error { return nil }
func (rh *redditHandler) Mention(m *reddit.Message) error          { return nil }

func (a *ActivityRedditListener) GetTargets(ctx context.Context) (RunActRequest, error) {
	r, err := http.NewRequest(
		http.MethodGet,
		os.Getenv("KAGGO_ENDPOINT")+"/notification/reddit/targets",
		nil)
	if err != nil {
		return RunActRequest{}, err
	}
	r.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return RunActRequest{}, err
	}
	if res.StatusCode != http.StatusOK {
		return RunActRequest{}, fmt.Errorf("bad response: %s", res.Status)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return RunActRequest{}, err
	}
	var ar RunActRequest
	err = json.Unmarshal(b, &ar)
	if err != nil {
		return RunActRequest{}, err
	}
	return ar, nil
}

func (a *ActivityRedditListener) Run(ctx context.Context, r RunActRequest) error {
	l := activity.GetLogger(ctx)

	cfg := reddit.BotConfig{
		Agent: os.Getenv("REDDIT_USER_AGENT"),
		App: reddit.App{
			ID:       os.Getenv("REDDIT_CLIENT_ID"),
			Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
			Username: os.Getenv("REDDIT_USERNAME"),
			Password: os.Getenv("REDDIT_PASSWORD"),
		},
		Rate: 5 * time.Second,
	}
	bot, err := reddit.NewBot(cfg)
	if err != nil {
		return err
	}
	lcfg := graw.Config{
		Subreddits: r.Subreddits,
		Users:      r.Users,
	}
	_, wait, err := graw.Run(&redditHandler{}, bot, lcfg)
	if err != nil {
		return fmt.Errorf("error running graw: %w", err)
	}
	l.Info("starting reddit listener", "subreddits", r.Subreddits, "users", r.Users)

	// continuously wait for graw updates, send any errors to the channel
	errC := make(chan error)
	go func() {
		for {
			if err := wait(); err != nil {
				errC <- err
			}
		}
	}()

	outer := true
	ticker := time.NewTicker(10 * time.Second)
	for outer {
		select {
		case <-ctx.Done():
			return context.Canceled
		case <-ticker.C:
			activity.RecordHeartbeat(ctx)
		case err := <-errC:
			l.Error("encountered graw error", "error", err.Error())
		}
	}
	return nil
}
