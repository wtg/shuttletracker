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

// casClient stores the local cas client and an instance of the database
type casClient struct {
	cas auth.AuthenticationService
	db  database.Database
}

// CreatecasClient creates an authentication service casClient using a cas url and database
func CreateCASClient(url *url.URL, db database.Database) (*casClient){
	client := gc.NewClient(&gc.Options{
		URL:   url,
		Store: nil,
	})

	cli := &casClient{
		cas: &auth.CAS{
			CAS: client,
		},
		db: db,
	}
	return cli
}

// InjectMocks allows mock interfaces to be used
func InjectMocks(cli auth.AuthenticationService, db database.Database) (*casClient){

	c:= &casClient{
		cas: cli,
		db: db,
	}
	return c
}

func (cli *casClient) logout(w http.ResponseWriter, r *http.Request) {
	cli.cas.Logout(w, r)
}

func (cli *casClient) casauth(next http.Handler) http.Handler {
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
