package temporal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/brojonat/kaggo/server/db/dbgen"
)

type DoRequestParam struct {
	Serial []byte `json:"serial"`
}
type DoRequestResult struct {
	StatusCode int    `json:"status_code"`
	Body       []byte `json:"body"`
}

type UploadResponseParam struct {
	ResponseKind string `json:"response_kind"`
	Serial       []byte `json:"serial"`
}
type UploadResponseResult struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ActivityRequester struct {
	DB *dbgen.Queries
}

func DoRequest(ctx context.Context, drp DoRequestParam) (*DoRequestResult, error) {
	buf := bufio.NewReader(bytes.NewReader(drp.Serial))
	r, err := http.ReadRequest(buf)
	if err != nil {
		return nil, fmt.Errorf("error deserializing request: %w", err)
	}

	// https://stackoverflow.com/questions/19595860/http-request-requesturi-field-when-making-request-in-go
	// RequestURI must not be set, but req.URL is incomplete, so parse a new one
	u, err := url.Parse("https://" + r.Host + r.RequestURI)
	if err != nil {
		return nil, fmt.Errorf("error parsing request URL: %s", r.RequestURI)
	}
	r.URL = u
	r.RequestURI = ""

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	res := DoRequestResult{
		StatusCode: resp.StatusCode,
		Body:       b,
	}
	return &res, nil
}

// UploadMetrics will handle the response from a get metrics
func UploadMetrics(ctx context.Context, kind string, b []byte) (*UploadResponseResult, error) {
	switch kind {
	case "internal.random":
		return handleInternalRandomMetrics(b)
	case "youtube.video":
		return handleYouTubeVideoMetrics(b)
	case "kaggle.notebook":
		return handleKaggleNotebookMetrics(b)
	case "kaggle.dataset":
		return handleKaggleDatasetMetrics(b)
	default:
		return nil, fmt.Errorf("unrecognized response kind: %s", kind)
	}
}
