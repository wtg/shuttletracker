package api

import (
	"gopkg.in/cas.v2"
	"net/http"
	"net/url"
)

type CasClient struct {
	cas    *cas.Client
	tstore *cas.MemoryStore
}

func (cli *CasClient) Create(url *url.URL, tickets *cas.MemoryStore) {
	client := cas.NewClient(&cas.Options{
		URL:   url,
		Store: nil,
	})
	cli.cas = client
	cli.tstore = tickets
}

func (cli *CasClient) logout(w http.ResponseWriter, r *http.Request) {
	cas.RedirectToLogout(w, r)
}

func (cli *CasClient) casauth(next http.Handler) http.Handler {
	return cli.cas.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")

		if !cas.IsAuthenticated(r) {
			cas.RedirectToLogin(w, r)
			// return

		} else {
		}
		next.ServeHTTP(w, r)

	})
}
