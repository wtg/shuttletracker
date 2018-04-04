package api

import (
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
	"gopkg.in/cas.v2"
	"net/http"
	"strings"
)

// casClient stores the local cas client and an instance of the database
type casClient struct {
	cas *cas.Client
	db  database.Database
}

func (cli *casClient) logout(w http.ResponseWriter, r *http.Request) {
	cas.RedirectToLogout(w, r)
}

func (cli *casClient) casauth(next http.Handler) http.Handler {
	return cli.cas.HandleFunc(func(w http.ResponseWriter, r *http.Request) {

		if !cas.IsAuthenticated(r) {
			_, err := w.Write([]byte("redirecting to cas;"))
			if err != nil {
				log.WithError(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			cas.RedirectToLogin(w, r)
		} else {
			auth, err := cli.db.UserExists(strings.ToLower(cas.Username(r)))
			if err != nil {
				log.WithError(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				auth = false
			}
			if auth {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "unauthenticated", 401)

		}

	})
}
