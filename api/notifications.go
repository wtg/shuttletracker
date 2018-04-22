package api

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
)

// AdminMessageHandler handles the retrieval of the current administrator message
func (api *API) AdminMessageHandler(w http.ResponseWriter, r *http.Request) {
	message, err := api.msg.Message()
	if err != nil {
		log.WithError(err).Error("unable to get message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, message)
}

// SetAdminMessage allows the user to set an alert message that will display to all users who visit the page
func (api *API) SetAdminMessage(w http.ResponseWriter, r *http.Request) {
	message := &shuttletracker.Message{}
	err := json.NewDecoder(r.Body).Decode(message)
	if err != nil {
		log.WithError(err).Error("unable to decode message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(message.Message) > 250 {
		http.Error(w, "message too long, must be fewer than 251 characters", 400)
		return
	}
	message.Message = template.HTMLEscapeString(message.Message)
	err = api.msg.SetMessage(message)
	if err != nil {
		log.WithError(err).Error("unable to update message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, "Success")
}
