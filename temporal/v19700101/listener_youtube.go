package temporal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.temporal.io/sdk/activity"
)

func (a *ActivityYouTubeListener) GetYouTubeChannelTargets(ctx context.Context) (YouTubeChannelSubActRequest, error) {
	r, err := http.NewRequest(
		http.MethodGet,
		os.Getenv("KAGGO_ENDPOINT")+"/notification/youtube/targets",
		nil)
	if err != nil {
		return YouTubeChannelSubActRequest{}, err
	}
	r.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return YouTubeChannelSubActRequest{}, err
	}
	if res.StatusCode != http.StatusOK {
		return YouTubeChannelSubActRequest{}, fmt.Errorf("bad response: %s", res.Status)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return YouTubeChannelSubActRequest{}, err
	}
	var ar YouTubeChannelSubActRequest
	err = json.Unmarshal(b, &ar)
	if err != nil {
		return YouTubeChannelSubActRequest{}, err
	}
	return ar, nil
}

func (a *ActivityYouTubeListener) Subscribe(ctx context.Context, ar YouTubeChannelSubActRequest) error {
	l := activity.GetLogger(ctx)
	errIDs := []string{}
	for _, id := range ar.ChannelIDs {
		if err := a.webSubSub(id); err != nil {
			l.Error("error subscribing to websub", "id", id, "error", err.Error())
			errIDs = append(errIDs, id)
		}
	}
	if len(errIDs) > 0 {
		return fmt.Errorf("could not subscribe to the following IDs: %s", errIDs)
	}
	return nil
}

func (a *ActivityYouTubeListener) webSubSub(id string) error {
	data := url.Values{}
	data.Set("hub.callback", fmt.Sprintf("%s/notification/youtube/websub", os.Getenv("KAGGO_ENDPOINT")))
	data.Set("hub.mode", "subscribe")
	data.Set("hub.topic", fmt.Sprintf("https://www.youtube.com/xml/feeds/videos.xml?channel_id=%s", id))
	r, err := http.NewRequest(
		http.MethodPost,
		"https://pubsubhubbub.appspot.com",
		strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("bad response (%d) and could not read body", res.StatusCode)
		}
		return fmt.Errorf("bad response (%d): %s", res.StatusCode, string(b))
	}
	return nil
}
