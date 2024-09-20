package server

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"

	"github.com/brojonat/kaggo/server/db/dbgen"
	kt "github.com/brojonat/kaggo/temporal/v19700101"
)

var errUnsupportedRequestKind = errors.New("unsupported request kind")

// Helper function that creates a request, serializes it, and computes the id
// from the hash of the bytes. This is handy for passing to various workflows
// called in this package.
func makeExternalRequest(q *dbgen.Queries, rk, id string, isMeta bool) (*http.Request, []byte, string, error) {
	// construct request by switching over RequestKind
	var err error
	var rwf *http.Request
	switch rk {
	case kt.RequestKindInternalRandom:
		rwf, err = makeExternalRequestInternalRandom()
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
	case kt.RequestKindRedditSubredditMonitor:
		rwf, err = makeExternalRequestRedditSubredditMonitor(id)
		if err != nil {
			return nil, nil, "", err
		}
	case kt.RequestKindRedditUser:
		rwf, err = makeExternalRequestRedditUser(id)
		if err != nil {
			return nil, nil, "", err
		}
	case kt.RequestKindRedditUserMonitor:
		rwf, err = makeExternalRequestRedditUserMonitor(id)
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
		if isMeta {
			rwf, err = makeExternalRequestTwitchStreamMeta(id)
		} else {
			rwf, err = makeExternalRequestTwitchStream(id)
		}
		if err != nil {
			return nil, nil, "", err
		}
	case kt.RequestKindTwitchUserPastDec:
		if isMeta {
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
	if isMeta {
		rid = fmt.Sprintf("%s %s metadata", rk, id)
	} else {
		rid = fmt.Sprintf("%s %s %x", rk, id, h.Sum(nil))
	}
	return rwf, serialReq, rid, nil
}

// NOTE: the following makeExternalRequestFoo functions will return a PROTOTYPE
// of a request to be made to an external resource. The implementation of these
// functions MUST NOT include dynamic values! The returned request may be
// serialized and a hash of the request bytes may be used as part of the
// schedule identifier. If any dynamic values are included and later on change,
// the hash of the resultant request will also change, and we won't be able to
// trivially prevent duplicate schedules from being created.

func makeExternalRequestInternalRandom() (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.kaggo.brojonat.com/internal/generate", nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func makeExternalRequestYouTubeVideo(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://youtube.googleapis.com/youtube/v3/videos", nil)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()
	q.Set("part", "snippet,contentDetails,statistics")
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
	return r, nil
}

func makeExternalRequestRedditSubreddit(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://oauth.reddit.com/r/%s/about.json", id), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func makeExternalRequestRedditSubredditMonitor(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://oauth.reddit.com/r/%s.json", id), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func makeExternalRequestRedditUser(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://oauth.reddit.com/user/%s/about.json", id), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func makeExternalRequestRedditUserMonitor(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://oauth.reddit.com/user/%s/submitted.json", id), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func makeExternalRequestTwitchClip(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/clips", nil)
	if err != nil {
		return nil, err
	}
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
	q := r.URL.Query()
	q.Add("id", id)
	r.URL.RawQuery = q.Encode()
	return r, nil
}

func makeExternalRequestTwitchStreamMeta(username string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return nil, err
	}
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
			user_id = md.Data.UserID
			break
		}
	}
	if user_id == "" {
		return nil, fmt.Errorf("error getting twitch user_id from metadata: no user_id found in metadata")
	}
	qs := r.URL.Query()
	qs.Add("user_id", user_id)
	r.URL.RawQuery = qs.Encode()
	return r, nil
}
