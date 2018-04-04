package cas

import (
	"net/http"
)

type Cas interface {
	Authenticated(request *http.Request) bool
	Logout(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Username(r *http.Request) string
	HandleFunc(func(http.ResponseWriter, *http.Request)) http.Handler
}
