package temporal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/brojonat/kaggo/server/api"
	"github.com/jmespath/go-jmespath"
)

// this is a convenience function for testing an internal random metric
func handleInternalRandomMetrics(b []byte) (*UploadResponseResult, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing internal response: %w", err)
	}

	// slug
	iface, err := jmespath.Search("id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting internal id: %w", err)
	}
	id := iface.(string)

	// value
	iface, err = jmespath.Search("value", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting internal value: %w", err)
	}
	value := iface.(float64)

	payload := api.InternalMetricPayload{
		ID:    id,
		Value: int(value),
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	r, err := http.NewRequest(http.MethodPost, "https://api.kaggo.brojonat.com/internal/metrics", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("error making request to upload metrics: %w", err)
	}
	r.Header.Add("Authorization", os.Getenv("AUTH_TOKEN")) // FIXME: inject headers?
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error doing request to upload metrics: %w", err)
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading metrics response body: %w", err)
	}
	var mr api.DefaultJSONResponse
	err = json.Unmarshal(b, &mr)
	if err != nil {
		return nil, fmt.Errorf("error parsing metrics response: %w", err)
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("bad response code uploading metrics: %d (%s): %s", res.StatusCode, http.StatusText(res.StatusCode), b)
	}
	urr := UploadResponseResult(mr)
	return &urr, nil
}

func handleKaggleNotebookMetrics(b []byte) (*UploadResponseResult, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing notebook response: %w", err)
	}

	// slug
	iface, err := jmespath.Search("notebook.slug.search.path", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting notebook slug: %w", err)
	}
	slug := iface.(string)

	// votes
	iface, err = jmespath.Search("notebook.votes.search.path", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting notebook votes: %w", err)
	}
	votes := iface.(float64)

	// downloads
	iface, err = jmespath.Search("notebook.downloads.search.path", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting notebook downloads: %w", err)
	}
	downloads := iface.(float64)

	// upload the metrics to the server
	payload := api.KaggleMetricPayload{
		Slug:         slug,
		Votes:        int(votes),
		Downloads:    int(downloads),
		SetVotes:     true,
		SetDownloads: true,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	r, err := http.NewRequest(http.MethodPost, "https://api.kaggo.brojonat.com/kaggle/notebook", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("error making request to upload metrics: %w", err)
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error doing request to upload metrics: %w", err)
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading metrics response body: %w", err)
	}
	var mr api.DefaultJSONResponse
	err = json.Unmarshal(b, &mr)
	if err != nil {
		return nil, fmt.Errorf("error parsing metrics response: %w", err)
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error uploading metrics: %d (%s): %s", res.StatusCode, http.StatusText(res.StatusCode), b)
	}
	urr := UploadResponseResult(mr)
	return &urr, nil
}

func handleKaggleDatasetMetrics(b []byte) (*UploadResponseResult, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing dataset response: %w", err)
	}

	// slug
	iface, err := jmespath.Search("dataset.slug.search.path", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting dataset slug: %w", err)
	}
	slug := iface.(string)

	// votes
	iface, err = jmespath.Search("dataset.votes.search.path", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting dataset votes: %w", err)
	}
	votes := iface.(float64)

	// downloads
	iface, err = jmespath.Search("dataset.downloads.search.path", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting dataset downloads: %w", err)
	}
	downloads := iface.(float64)

	// upload the metrics to the server
	payload := api.KaggleMetricPayload{
		Slug:         slug,
		Votes:        int(votes),
		Downloads:    int(downloads),
		SetVotes:     true,
		SetDownloads: true,
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	r, err := http.NewRequest(http.MethodPost, "https://api.kaggo.brojonat.com/kaggle/dataset", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("error making request to upload metrics: %w", err)
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error doing request to upload metrics: %w", err)
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading metrics response body: %w", err)
	}
	var mr api.DefaultJSONResponse
	err = json.Unmarshal(b, &mr)
	if err != nil {
		return nil, fmt.Errorf("error parsing metrics response: %w", err)
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error uploading metrics: %d (%s): %s", res.StatusCode, http.StatusText(res.StatusCode), b)
	}
	urr := UploadResponseResult(mr)
	return &urr, nil
}

// this is a convenience function for testing an internal random metric
func handleYouTubeVideoMetrics(b []byte) (*UploadResponseResult, error) {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error deserializing internal response: %w", err)
	}

	// id
	iface, err := jmespath.Search("id", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting id: %w", err)
	}
	id := iface.(string)

	// slug
	iface, err = jmespath.Search("slug", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting slug: %w", err)
	}
	slug := iface.(string)

	// views
	iface, err = jmespath.Search("views", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting views: %w", err)
	}
	views := iface.(float64)

	// comments
	iface, err = jmespath.Search("comments", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting comments: %w", err)
	}
	comments := iface.(float64)

	// likes
	iface, err = jmespath.Search("likes", data)
	if err != nil {
		return nil, fmt.Errorf("error extracting likes: %w", err)
	}
	likes := iface.(float64)

	payload := api.YouTubeVideoMetricPayload{
		ID:          id,
		Slug:        slug,
		SetViews:    true,
		Views:       int(views),
		SetComments: true,
		Comments:    int(comments),
		SetLikes:    true,
		Likes:       int(likes),
	}
	b, err = json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing upload data: %w", err)
	}
	r, err := http.NewRequest(http.MethodPost, "https://api.kaggo.brojonat.com/youtube/video", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("error making request to upload metrics: %w", err)
	}
	r.Header.Add("Authorization", os.Getenv("AUTH_TOKEN")) // FIXME: inject headers?
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error doing request to upload metrics: %w", err)
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading metrics response body: %w", err)
	}
	var mr api.DefaultJSONResponse
	err = json.Unmarshal(b, &mr)
	if err != nil {
		return nil, fmt.Errorf("error parsing metrics response: %w", err)
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("bad response code uploading metrics: %d (%s): %s", res.StatusCode, http.StatusText(res.StatusCode), b)
	}
	urr := UploadResponseResult(mr)
	return &urr, nil
}
