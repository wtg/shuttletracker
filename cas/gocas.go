package cas

import (
	"gopkg.in/cas.v2"
	"net/http"
)

type Gocas struct {
	Cas *cas.Client
}

// Authenticated returns true if the request has been authenticated with cas, false otherwise
func (c *Gocas) Authenticated(request *http.Request) bool {
	return cas.IsAuthenticated(request)
}

// Logout redirects the user to logout
func (c *Gocas) Logout(w http.ResponseWriter, r *http.Request) {
	cas.RedirectToLogout(w, r)
}

// Login redirects the user to login
func (c *Gocas) Login(w http.ResponseWriter, r *http.Request) {
	cas.RedirectToLogin(w, r)

}

// Username returns the username of the request
func (c *Gocas) Username(r *http.Request) string {
	return cas.Username(r)
}

func (c *Gocas) HandleFunc(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return c.Cas.HandleFunc(f)
}
