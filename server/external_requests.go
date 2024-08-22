package server

import (
	"fmt"
	"net/http"
	"os"
)

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
