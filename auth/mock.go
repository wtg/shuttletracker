package auth

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

// Mock implements the auth interface.
type Mock struct {
	mock.Mock
}

// Authenticated returns the mock response to the server
func (auth *Mock) Authenticated(request *http.Request) bool {
	return true
}

// Logout writes logout to the ResponseWriter
func (auth *Mock) Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("logout"))

}

// Login writes login to the ResponseWriter
func (auth *Mock) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))

}

// Username returns the mock response to the server
func (auth *Mock) Username(request *http.Request) string {
	return "lyonj4"
}

func (auth *Mock) HandleFunc(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(f)
}
