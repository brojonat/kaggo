package server

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/brojonat/kaggo/server/db/dbgen"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
)

var errUnsupportedRequestKind = errors.New("unsupported request kind")

// Helper function that creates a request, serializes it, and computes the id from the hash
// of the bytes. This is handy for passing to various workflows called in this package.
func makeExternalRequest(q *dbgen.Queries, rk, id string, meta bool) (*http.Request, []byte, string, error) {
	// construct request by switching over RequestKind
	var err error
	var rwf *http.Request
	switch rk {
	case kt.RequestKindInternalRandom:
		rwf, err = makeExternalRequestInternalRandom()
		if err != nil {
			return nil, nil, "", err
		}

	case kt.RequestKindYouTubeVideo:
		rwf, err = makeExternalRequestYouTubeVideo(id)
		if err != nil {
			return nil, nil, "", err
		}

	case kt.RequestKindYouTubeChannel:
		rwf, err = makeExternalRequestYouTubeChannel(id)
		if err != nil {
			return nil, nil, "", err
		}

	case kt.RequestKindKaggleNotebook:
		rwf, err = makeExternalRequestKaggleNotebook(id)
		if err != nil {
			return nil, nil, "", err
		}

	case kt.RequestKindKaggleDataset:
		rwf, err = makeExternalRequestKaggleDataset(id)
		if err != nil {
			return nil, nil, "", err
		}

	case kt.RequestKindRedditPost:
		rwf, err = makeExternalRequestRedditPost(id)
		if err != nil {
			return nil, nil, "", err
		}

	case kt.RequestKindRedditComment:
		rwf, err = makeExternalRequestRedditComment(id)
		if err != nil {
			return nil, nil, "", err
		}

	case kt.RequestKindRedditSubreddit:
		rwf, err = makeExternalRequestRedditSubreddit(id)
		if err != nil {
			return nil, nil, "", err
		}

	case kt.RequestKindTwitchClip:
		rwf, err = makeExternalRequestTwitchClip(id)
		if err != nil {
			return nil, nil, "", err
		}
	case kt.RequestKindTwitchVideo:
		rwf, err = makeExternalRequestTwitchVideo(id)
		if err != nil {
			return nil, nil, "", err
		}
	case kt.RequestKindTwitchStream:
		if meta {
			rwf, err = makeExternalRequestTwitchStreamMeta(id)
		} else {
			rwf, err = makeExternalRequestTwitchStream(id)
		}
		if err != nil {
			return nil, nil, "", err
		}
	case kt.RequestKindTwitchUserPastDec:
		if meta {
			rwf, err = makeExternalRequestTwitchUserPastDecMeta(id)
		} else {
			rwf, err = makeExternalRequestTwitchUserPastDec(q, id)
		}
		if err != nil {
			return nil, nil, "", err
		}

	default:
		return nil, nil, "", errUnsupportedRequestKind
	}

	// serialize the request
	buf := &bytes.Buffer{}
	err = rwf.Write(buf)
	if err != nil {
		return nil, nil, "", err
	}
	serialReq := buf.Bytes()
	h := md5.New()
	_, err = h.Write(serialReq)
	if err != nil {
		return nil, nil, "", err
	}

	// For polling requests, the identifier is the request kind, identifier, and
	// hash of the request. For metadata requests, the hash is omitted in favor
	// of a string that simply indicates "metadata".
	var rid string
	if meta {
		rid = fmt.Sprintf("%s %s metadata", rk, id)
	} else {
		rid = fmt.Sprintf("%s %s %x", rk, id, h.Sum(nil))
	}
	return rwf, serialReq, rid, nil
}

func makeExternalRequestInternalRandom() (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.kaggo.brojonat.com/internal/generate", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	return r, nil
}

func makeExternalRequestYouTubeVideo(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://youtube.googleapis.com/youtube/v3/videos", nil)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()
	q.Set("part", "snippet,contentDetails,statistics")
	q.Set("key", os.Getenv("YOUTUBE_API_KEY"))
	q.Set("id", id)
	r.URL.RawQuery = q.Encode()
	return r, nil
}

func makeExternalRequestYouTubeChannel(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://youtube.googleapis.com/youtube/v3/channels", nil)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()
	q.Set("part", "snippet,contentDetails,statistics")
	q.Set("key", os.Getenv("YOUTUBE_API_KEY"))
	q.Set("id", id)
	r.URL.RawQuery = q.Encode()
	return r, nil
}

func makeExternalRequestKaggleNotebook(id string) (*http.Request, error) {
	// https://github.com/Kaggle/kaggle-api/blob/48d0433575cac8dd20cf7557c5d749987f5c14a2/kaggle/api/kaggle_api.py#L3052
	r, err := http.NewRequest(http.MethodGet, "https://www.kaggle.com/api/v1/kernels/list", nil)
	if err != nil {
		return nil, err
	}
	// filter to notebooks only and search using the supplied ref
	q := r.URL.Query()
	q.Set("search", id)
	r.URL.RawQuery = q.Encode()
	// basic auth
	r.SetBasicAuth(os.Getenv("KAGGLE_USERNAME"), os.Getenv("KAGGLE_API_KEY"))
	return r, nil
}

func makeExternalRequestKaggleDataset(id string) (*http.Request, error) {
	// https: //github.com/Kaggle/kaggle-api/blob/48d0433575cac8dd20cf7557c5d749987f5c14a2/kaggle/api/kaggle_api.py#L1731
	r, err := http.NewRequest(http.MethodGet, "https://www.kaggle.com/api/v1/datasets/list", nil)
	if err != nil {
		return nil, err
	}
	// search using the supplied ref
	q := r.URL.Query()
	q.Set("search", id)
	r.URL.RawQuery = q.Encode()
	// basic auth
	r.SetBasicAuth(os.Getenv("KAGGLE_USERNAME"), os.Getenv("KAGGLE_API_KEY"))
	return r, nil
}

func makeExternalRequestRedditPost(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://oauth.reddit.com/api/info.json", nil)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()
	q.Set("id", fmt.Sprintf("t3_%s", id))
	r.URL.RawQuery = q.Encode()
	r.Header.Add("User-Agent", "Debian:github.com/brojonat/kaggo/worker:v0.0.1 (by /u/GraearG)")
	return r, nil
}

func makeExternalRequestRedditComment(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://oauth.reddit.com/api/info.json", nil)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()
	q.Set("id", fmt.Sprintf("t1_%s", id))
	r.URL.RawQuery = q.Encode()
	r.Header.Add("User-Agent", "Debian:github.com/brojonat/kaggo/worker:v0.0.1 (by /u/GraearG)")
	return r, nil
}

func makeExternalRequestRedditSubreddit(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://oauth.reddit.com/r/%s/about.json", id), nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("User-Agent", "Debian:github.com/brojonat/kaggo/worker:v0.0.1 (by /u/GraearG)")
	return r, nil
}

func makeExternalRequestTwitchClip(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/clips", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))
	q := r.URL.Query()
	q.Add("id", id)
	r.URL.RawQuery = q.Encode()
	return r, nil
}

func makeExternalRequestTwitchVideo(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/videos", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))
	q := r.URL.Query()
	q.Add("id", id)
	r.URL.RawQuery = q.Encode()
	return r, nil
}

func makeExternalRequestTwitchStreamMetadata(username string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/streams", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))
	q := r.URL.Query()
	q.Add("user_login", username)
	q.Add("sort", "time")
	r.URL.RawQuery = q.Encode()
	return r, nil
}

func makeExternalRequestTwitchStreamMeta(username string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))
	q := r.URL.Query()
	q.Add("login", username)
	q.Add("sort", "time")
	r.URL.RawQuery = q.Encode()
	return r, nil
}

func makeExternalRequestTwitchStream(username string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/streams", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))
	q := r.URL.Query()
	q.Add("user_login", username)
	q.Add("sort", "time")
	r.URL.RawQuery = q.Encode()
	return r, nil
}

func makeExternalRequestTwitchUserPastDecMeta(username string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))
	qs := r.URL.Query()
	qs.Add("login", username)
	r.URL.RawQuery = qs.Encode()
	return r, nil
}

func makeExternalRequestTwitchUserPastDec(q *dbgen.Queries, username string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/videos", nil)
	if err != nil {
		return nil, err
	}

	// This assumes that we've already run the metadata workflow for this
	// request. This is necessary because we need the Twitch user_id to fetch a
	// user's recent videos. Our choices are to query the Twitch API here, or
	// leverage the fact that the application has already run the metadata
	// workflow and that the corresponding entry exists in the DB under the
	// supplied username. If a caller is pathological, they can find a way to
	// skip the metadata workflow, but then this will simply return an error.
	mds, err := q.GetMetadataByIDs(context.Background(), []string{username})
	if err != nil {
		return nil, fmt.Errorf("error getting twitch user_id from metadata: %w", err)
	}
	if len(mds) == 0 {
		return nil, fmt.Errorf("error getting twitch user_id from metadata: no result rows")
	}

	var user_id string
	for _, md := range mds {
		if md.RequestKind == kt.RequestKindTwitchUserPastDec {
			user_id = md.Data.ID
			break
		}
	}
	if user_id == "" {
		return nil, fmt.Errorf("error getting twitch user_id from metadata: no user_id found in metadata")
	}
	r.Header.Add("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))
	qs := r.URL.Query()
	qs.Add("user_id", user_id)
	r.URL.RawQuery = qs.Encode()
	return r, nil
}
