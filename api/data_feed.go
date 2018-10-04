package api

import (
	"net/http"

	"github.com/wtg/shuttletracker/log"
)

// DataFeedHandler returns the latest successful response that the Updater received
// from the iTRAK data feed.
func (api *API) DataFeedHandler(w http.ResponseWriter, r *http.Request) {
	dfresp := api.updater.GetLastResponse()
	if dfresp == nil {
		http.Error(w, "Last data feed response does not exist", http.StatusNotFound)
		return
	}
	_, err := w.Write(dfresp.Body)
	if err != nil {
		log.WithError(err).Error("unable to write")
	}
}
