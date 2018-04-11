package api

import (
	"github.com/wtg/shuttletracker/auth"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
	"net/url"

	gc "gopkg.in/cas.v2"

	"net/http"
	"strings"
)

// CasClient stores the local cas client and an instance of the database
type CasClient struct {
	cas auth.AuthenticationService
	db  database.Database
}

func CreateCasClient(url *url.URL, db database.Database) (*CasClient){
	client := gc.NewClient(&gc.Options{
		URL:   url,
		Store: nil,
	})

	cli := &CasClient{
		cas: &auth.CAS{
			CAS: client,
		},
		db: db,
	}
	return cli
}

func Clone(cli auth.AuthenticationService, db database.Database) (*CasClient){

	c:= &CasClient{
		cas: cli,
		db: db,
	}
	return c
}

func (cli *CasClient) logout(w http.ResponseWriter, r *http.Request) {
	cli.cas.Logout(w, r)
}

func (cli *CasClient) casauth(next http.Handler) http.Handler {
	return cli.cas.HandleFunc(func(w http.ResponseWriter, r *http.Request) {

		if !cli.cas.Authenticated(r) {
			_, err := w.Write([]byte("redirecting to cas;"))
			if err != nil {
				log.WithError(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			cli.cas.Login(w, r)
		} else {
			auth, err := cli.db.UserExists(strings.ToLower(cli.cas.Username(r)))
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
