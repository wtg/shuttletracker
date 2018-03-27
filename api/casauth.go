package api

import (
	"github.com/wtg/shuttletracker/database"
	"gopkg.in/cas.v2"
	"net/http"
	"net/url"
	"strings"
)

type CasClient struct {
	cas *cas.Client
	db  database.Database
}

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
			users, _ := cli.db.GetUsers()
			valid := false
			for _, u := range users {
				if u.Name == strings.ToLower(cas.Username(r)) {
					valid = true
				}
			}
			if !valid {
				cas.RedirectToLogout(w, r)
				return
			}
			next.ServeHTTP(w, r)
		}

	})
}
