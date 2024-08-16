package temporal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func (a *ActivityRequester) ensureValidRedditToken(minDur time.Duration) error {
	// reddit@reddit-VirtualBox:~$ curl -X POST -d 'grant_type=password&username=reddit_bot&password=snoo' --user 'p-jcoLKBynTLew:gko_LXELoV07ZBNUXrvWZfzE3aI' https://www.reddit.com/api/v1/access_token
	// {
	// 	"access_token": "J1qK1c18UUGJFAzz9xnH56584l4",
	// 	"expires_in": 3600,
	// 	"scope": "*",
	// 	"token_type": "bearer"
	// }

	// short circuit early if the token doesn't need to be refreshed
	if time.Until(a.RedditAuthTokenExp) > minDur {
		return nil
	}

	// otherwise hit the reddit API for a new token
	formData := url.Values{
		"grant_type": {"password"},
		"username":   {os.Getenv("REDDIT_USERNAME")},
		"password":   {os.Getenv("REDDIT_PASSWORD")},
	}
	r, err := http.NewRequest(http.MethodPost, "https://www.reddit.com/api/v1/access_token", strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	r.SetBasicAuth(os.Getenv("REDDIT_CLIENT_ID"), os.Getenv("REDDIT_CLIENT_SECRET"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &body)
	if err != nil {
		return err
	}

	a.RedditAuthToken = body.AccessToken
	dur := time.Duration(body.ExpiresIn * int(time.Second))
	a.RedditAuthTokenExp = time.Now().Add(dur)
	return nil
}
