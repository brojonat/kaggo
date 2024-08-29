package temporal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func (a *ActivityRequester) ensureValidTwitchToken(minDur time.Duration) error {
	// {
	// 	"access_token": "jostpf5q0uzmxmkba9iyug38kjtgh",
	// 	"expires_in": 5011271,
	// 	"token_type": "bearer"
	//   }
	// short circuit early if the token doesn't need to be refreshed
	if time.Until(a.TwitchAuthTokenExp) > minDur {
		return nil
	}

	// otherwise hit the reddit API for a new token
	formData := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {os.Getenv("TWITCH_CLIENT_ID")},
		"client_secret": {os.Getenv("TWITCH_CLIENT_SECRET")},
	}
	r, err := http.NewRequest(http.MethodPost, "https://id.twitch.tv/oauth2/token", strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad response %d code for getting twitch auth token: %s", resp.StatusCode, string(b))
	}
	err = json.Unmarshal(b, &body)
	if err != nil {
		return err
	}

	a.TwitchAuthToken = body.AccessToken
	dur := time.Duration(body.ExpiresIn * int(time.Second))
	a.TwitchAuthTokenExp = time.Now().Add(dur)
	return nil
}
