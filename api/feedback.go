package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wtg/shuttletracker"
)

// FeedbackHandler finds all forms in the database
func (api *API) FeedbackHandler(w http.ResponseWriter, r *http.Request) {
	forms, err := api.fdb.Forms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, forms)
}

// FeedbackCreateHandler adds a new form to the database
func (api *API) FeedbackCreateHandler(w http.ResponseWriter, r *http.Request) {
	form := &shuttletracker.Form{}
	err := json.NewDecoder(r.Body).Decode(form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// originally put create form here, but feedback forms dont have manual create

}

// FormsEditHandler (idk if needed)

// FeedbackDeleteHandler deletes a form from database
func (api *API) FeedbackDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.fdb.DeleteForm(id)
	if err != nil {
		if err == shuttletracker.ErrVehicleNotFound {
			http.Error(w, "Form not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// FeedbackEditHandler Only handles editing enabled flag for now
func (api *API) FeedbackEditHandler(w http.ResponseWriter, r *http.Request) {
	form := &shuttletracker.Form{}
	err := json.NewDecoder(r.Body).Decode(form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rd := form.Read
	form, err = api.fdb.Form(form.ID)
	form.Read = rd
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.fdb.EditForm(form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
