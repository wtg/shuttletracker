package auth

import (
	"net/http"
)

// AuthenticationService interface represents an authentication service for login
type AuthenticationService interface {
	Authenticated(r *http.Request) bool
	Logout(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Username(r *http.Request) string
	HandleFunc(func(http.ResponseWriter, *http.Request)) http.Handler
}
