package temporal

// workflows

type DoMetadataRequestWFRequest struct {
	RequestKind string `json:"request_kind"`
	Serial      []byte `json:"serial"`
}

type DoPollingRequestWFRequest struct {
	RequestKind string `json:"request_kind"`
	Serial      []byte `json:"serial"`
}

// activities

type DoRequestActRequest struct {
	RequestKind string `json:"request_kind"`
	Serial      []byte `json:"serial"`
}
type DoRequestActResult struct {
	RequestKind string `json:"request_kind"`
	StatusCode  int    `json:"status_code"`
	Body        []byte `json:"body"`
}

type UploadMetadataActRequest struct {
	RequestKind string `json:"request_kind"`
	StatusCode  int    `json:"status_code"`
	Body        []byte `json:"body"`
}

type UploadMetricsActRequest struct {
	RequestKind string `json:"request_kind"`
	StatusCode  int    `json:"status_code"`
	Body        []byte `json:"body"`
}
