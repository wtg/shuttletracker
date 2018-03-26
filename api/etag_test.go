package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestETag(t *testing.T) {
	r := chi.NewRouter()
	r.Use(etag)

	ts := httptest.NewServer(r)
	defer ts.Close()

	type testCase struct {
		clientETag string
		serverETag string
		statusCode int
		body       string
	}
	cases := []testCase{
		{
			clientETag: "",
			serverETag: "6e71b3cac15d32fe2d36c270887df9479c25c640",
			statusCode: 200,
			body:       "hello there",
		},
		{
			clientETag: "6e71b3cac15d32fe2d36c270887df9479c25c640",
			serverETag: "6e71b3cac15d32fe2d36c270887df9479c25c640",
			statusCode: 304,
			body:       "hello there",
		},
	}

	for _, c := range cases {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(c.body))
		})

		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			t.Errorf("unable to create HTTP request: %s", err)
			continue
		}
		if len(c.clientETag) != 0 {
			req.Header.Add("If-None-Match", c.clientETag)
		}

		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("unable to get URL: %s", err)
			continue
		}

		bodyBytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Errorf("unable to read: %s", err)
			continue
		}
		body := string(bodyBytes)

		// server should return empty body if ETag matches
		if c.statusCode == 304 {
			if len(body) != 0 {
				t.Errorf("expected empty body. got [%s]", body)
			}
		} else {
			if body != c.body {
				t.Errorf("bodies not equal. got [%s], expected [%s]", body, c.body)
			}
		}

		etag := res.Header.Get("ETag")

		if etag != c.serverETag {
			t.Errorf("ETags not equal. got [%s], expected [%s]", etag, c.serverETag)
		}
	}
}
