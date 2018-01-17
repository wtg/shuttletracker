package api

import (
	"github.com/wtg/shuttletracker/model"
	"net/http"
)

func (api *API) AdminMessageHandler(w http.ResponseWriter, r *http.Request) {
	message := model.AdminMessage{}
	message.Type = "alert"
	message.Message = "This is a test message"
	message.Display = true
	WriteJSON(w, message)
}
