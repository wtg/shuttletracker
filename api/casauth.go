package api

import (
	"bytes"
	"context"
	"github.com/wtg/shuttletracker/log"
	"gopkg.in/cas.v2"
	"io"
	"net/http"
	"net/url"
)

type casResponseWriter struct {
	http.ResponseWriter
	buf bytes.Buffer
	w   io.Writer
	ctx context.Context
}

type CasClient struct {
	cas    *cas.Client
	tstore *cas.MemoryStore
}

func (e *casResponseWriter) Write(p []byte) (int, error) {
	return e.w.Write(p)
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

		ew := &casResponseWriter{
			ResponseWriter: w,
			buf:            bytes.Buffer{},
		}

		ew.w = io.Writer(&ew.buf)
		log.Debugf("here")
		if !cas.IsAuthenticated(r) {
			cas.RedirectToLogin(w, r)
			// return
		} else {
			next.ServeHTTP(ew, r)
			ew.buf.WriteTo(w)

		}

	})
}
