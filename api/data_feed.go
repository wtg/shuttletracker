package api

import (
	"net/http"
)

func (api *API) DataFeedHandler(w http.ResponseWriter, r *http.Request) {
	dfresp := api.updater.GetLastResponse()
	if dfresp == nil {
		http.Error(w, "Last data feed response does not exist", http.StatusNotFound)
		return
	}
	w.Write(dfresp.Body)
}
