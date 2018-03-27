package api

import (
	"github.com/wtg/shuttletracker/database"
	"gopkg.in/cas.v2"
	"net/http"
	"net/url"
	"strings"
)

// CasClient stores the local cas client and an instance of the database
type CasClient struct {
	cas *cas.Client
	db  database.Database
}

// Create creates a new CasClient from a casurl and a database
func (cli *CasClient) Create(url *url.URL, db database.Database) {
	client := cas.NewClient(&cas.Options{
		URL:   url,
		Store: nil,
	})
	cli.cas = client
	cli.db = db
}

func (cli *CasClient) logout(w http.ResponseWriter, r *http.Request) {
	cas.RedirectToLogout(w, r)
}

func (cli *CasClient) casauth(next http.Handler) http.Handler {
	return cli.cas.HandleFunc(func(w http.ResponseWriter, r *http.Request) {

		if !cas.IsAuthenticated(r) {
			cas.RedirectToLogin(w, r)
		} else {
			if cli.db.UserExists(strings.ToLower(cas.Username(r))) {
				next.ServeHTTP(w, r)
				return
			}
			cas.RedirectToLogout(w, r)
		}

	})
}
