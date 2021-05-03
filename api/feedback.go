package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wtg/shuttletracker"
)

// FeedbackAdminHandler gets the feedback message with admin=true
func (api *API) FeedbackAdminHandler(w http.ResponseWriter, r *http.Request) {
	form := api.fdb.GetAdminForm()
	WriteJSON(w, form)
}

// FeedbackHandler finds all forms in the database
func (api *API) FeedbackHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		forms, err := api.fdb.GetForms()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		WriteJSON(w, forms)
	} else {
		form, err := api.fdb.GetForm(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		WriteJSON(w, form)
	}
}

// FeedbackCreateHandler adds a new form to the database
func (api *API) FeedbackCreateHandler(w http.ResponseWriter, r *http.Request) {
	form := &shuttletracker.Form{}
	err := json.NewDecoder(r.Body).Decode(form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.fdb.CreateForm(form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FeedbackDeleteHandler deletes a form from database
func (api *API) FeedbackDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.fdb.DeleteForm(id)
	if err != nil {
		if err == shuttletracker.ErrFormNotFound {
			http.Error(w, "Form not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
