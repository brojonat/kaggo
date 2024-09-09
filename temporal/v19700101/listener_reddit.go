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

func (a *ActivityRedditListener) GetRedditUserTargets(ctx context.Context) (RedditSubActRequest, error) {
	r, err := http.NewRequest(
		http.MethodGet,
		os.Getenv("KAGGO_ENDPOINT")+"/notification/reddit/targets",
		nil)
	if err != nil {
		return RedditSubActRequest{}, err
	}
	r.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return RedditSubActRequest{}, err
	}
	if res.StatusCode != http.StatusOK {
		return RedditSubActRequest{}, fmt.Errorf("bad response: %s", res.Status)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return RedditSubActRequest{}, err
	}
	var ar RedditSubActRequest
	err = json.Unmarshal(b, &ar)
	if err != nil {
		return RedditSubActRequest{}, err
	}
	return ar, nil
}

func (a *ActivityRedditListener) Run(ctx context.Context, r RedditSubActRequest) error {
	l := activity.GetLogger(ctx)

	cfg := reddit.BotConfig{
		Agent: os.Getenv("REDDIT_LISTENER_USER_AGENT"),
		App: reddit.App{
			ID:       os.Getenv("REDDIT_LISTENER_CLIENT_ID"),
			Secret:   os.Getenv("REDDIT_LISTENER_CLIENT_SECRET"),
			Username: os.Getenv("REDDIT_LISTENER_USERNAME"),
			Password: os.Getenv("REDDIT_LISTENER_PASSWORD"),
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

	// need to send heartbeats while setting up since this can take a bit
	ticker := time.NewTicker(10 * time.Second)
	errC := make(chan error)
	var wait func() error
	go func() {
		tstart := time.Now()
		_, wait, err = graw.Run(&redditHandler{}, bot, lcfg)
		if err != nil {
			errC <- err
		}
		l.Info(
			"started reddit listener",
			"setup_duration", time.Since(tstart).String(),
			"n_subreddits", len(r.Subreddits),
			"n_users", len(r.Users),
		)
		errC <- nil
	}()
	doLoop := true
	for doLoop {
		select {
		case <-ticker.C:
			activity.RecordHeartbeat(ctx)
		case err := <-errC:
			if err != nil {
				return fmt.Errorf("error starting graw run: %w", err)
			}
			doLoop = false
		}
	}

	// Continuously wait for graw updates, send any errors to the channel.
	// FIXME: yes, I'm not cleaning this goroutine up as I should but the should
	// get killed with some regularity anyway so I don't anticipate it will be a
	// problem.
	go func() {
		for {
			if err := wait(); err != nil {
				errC <- err
			}
		}
	}()

	doLoop = true
	for doLoop {
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
