package api

import (
	"net/http"
)

// MessageGetHandler Returns the admin alert message when it is requested
func (api *API) MessageGetHandler(w http.ResponseWriter, r *http.Request) {
  stringThing := "Admin Alert Example:This is an example admin alert, it will be customizable"
  message := []byte(stringThing)
  _ = message
	WriteJSON(w, stringThing)
}
