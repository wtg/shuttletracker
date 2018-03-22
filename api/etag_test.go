package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
)

func TestETag(t *testing.T) {
	r := chi.NewRouter()

	r.Use(etag)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello there"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}

	msg, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(msg, []byte("hello there")) {
		t.Error("message not equal")
	}

	if res.Header.Get("ETag") != "6e71b3cac15d32fe2d36c270887df9479c25c640" {
		t.Errorf("ETag not equal. Got %s", res.Header.Get("ETag"))
	}
}
