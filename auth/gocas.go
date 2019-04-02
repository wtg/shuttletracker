package auth

import (
	"gopkg.in/cas.v2"
	"net/http"
)

// CAS is an implementation of the cas interface using the go-cas package
type CAS struct {
	CAS *cas.Client
}

// Authenticated returns true if the request has been authenticated with cas, false otherwise
func (c *CAS) Authenticated(request *http.Request) bool {
	return cas.IsAuthenticated(request)
}

// Logout redirects the user to logout
func (c *CAS) Logout(w http.ResponseWriter, r *http.Request) {
	cas.RedirectToLogout(w, r)
}

// Login redirects the user to login
func (c *CAS) Login(w http.ResponseWriter, r *http.Request) {
	cas.RedirectToLogin(w, r)

}

// Username returns the username of the request
func (c *CAS) Username(r *http.Request) string {
	return cas.Username(r)
}

//HandleFunc acts as an http handler for CAS
func (c *CAS) HandleFunc(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return c.CAS.HandleFunc(f)
}
