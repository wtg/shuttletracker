package shuttletracker

import (
	"net/http"
)

// DataFeedResponse contains information from the iTRAK data feed.
type DataFeedResponse struct {
	Body       []byte
	StatusCode int
	Headers    http.Header
}

// UpdaterService is an interface for interacting with vehicle location updates.
type UpdaterService interface {
	GetLastResponse() *DataFeedResponse
}
