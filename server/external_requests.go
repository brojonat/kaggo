package server

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"os"

	kt "github.com/brojonat/kaggo/temporal/v19700101"
)

var errUnsupportedRequestKind = errors.New("unsupported request kind")

// Helper function that creates a request, serializes it, and computes the id from the hash
// of the bytes. This is handy for passing to various workflows called in this package.
func makeExternalRequest(rk, id string) (*http.Request, []byte, string, error) {
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

	// identifier is the request kind, identifier, and hash of the request
	rid := fmt.Sprintf("%s %s %x", rk, id, h.Sum(nil))

	return rwf, serialReq, rid, nil
}

func makeExternalRequestInternalRandom() (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://api.kaggo.brojonat.com/internal/generate", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	r.Header.Add("Accept", "application/json")
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
	r.Header.Add("Accept", "application/json")
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
	r.Header.Add("Accept", "application/json")
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
	r.Header.Add("Accept", "application/json")
	r.SetBasicAuth(os.Getenv("KAGGLE_USERNAME"), os.Getenv("KAGGLE_API_KEY"))
	return r, nil
}

func makeExternalRequestRedditPost(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://reddit.com/api/info.json", nil)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()
	q.Set("id", fmt.Sprintf("t3_%s", id))
	r.URL.RawQuery = q.Encode()
	r.Header.Add("Accept", "application/json")
	r.Header.Add("User-Agent", "Debian:github.com/brojonat/kaggo/worker:v0.0.1 (by /u/GreaerG)")
	return r, nil
}

func makeExternalRequestRedditComment(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "https://reddit.com/api/info.json", nil)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()
	q.Set("id", fmt.Sprintf("t1_%s", id))
	r.URL.RawQuery = q.Encode()
	r.Header.Add("Accept", "application/json")
	r.Header.Add("User-Agent", "Debian:github.com/brojonat/kaggo/worker:v0.0.1 (by /u/GreaerG)")
	return r, nil
}

func makeExternalRequestRedditSubreddit(id string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://reddit.com/r/%s/about.json", id), nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Accept", "application/json")
	r.Header.Add("User-Agent", "Debian:github.com/brojonat/kaggo/worker:v0.0.1 (by /u/GreaerG)")
	return r, nil
}
