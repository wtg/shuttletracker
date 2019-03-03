package api

import (
	"net/http"
	"net/http/httputil"
	"strings"
)

func (api *API) createOSRMProxy(prefix string) http.Handler {
	proxy := &httputil.ReverseProxy{
		Director: func (r *http.Request) {
			r.Host = api.osrmUpstreamURL.Host
			r.URL.Host = api.osrmUpstreamURL.Host
			r.URL.Scheme = api.osrmUpstreamURL.Scheme
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		},
	}
	return proxy
}
