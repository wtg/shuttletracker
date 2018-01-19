package api

import (
	"github.com/wtg/shuttletracker/model"
	"net/http"
	"encoding/json"
)

// AdminMessageHandler handles the retrieval of the current administrator message
func (api *API) AdminMessageHandler(w http.ResponseWriter, r *http.Request) {
	message, err := api.db.GetCurrentMessage()

	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
	}
	WriteJSON(w, message)
}

// SetAdminMessage allows the user to set an alert message that will display to all users who visit the page
func (api *API) SetAdminMessage(w http.ResponseWriter, r *http.Request){
	message := model.AdminMessage{}
	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return;
	}

	err = api.db.SetMessage(&message)

	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
	}
}
