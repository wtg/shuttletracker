package api

import (
	"github.com/go-chi/chi"
	"github.com/wtg/shuttletracker/database"

	"gopkg.in/cas.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCasUnauthenticated(t *testing.T) {
	url, _ := url.Parse("https://cas.example.com/")
	client := cas.NewClient(&cas.Options{
		URL: url,
	})
	db := &database.Mock{}
	httpcli := http.Client{}
	cli := casClient{
		cas: client,
		db:  db,
	}

	r := chi.NewRouter()
	r.Use(cli.casauth)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Errorf("Error creating http request")
	}
	resp, err := httpcli.Do(req)
	if err != nil {
		t.Errorf("Error performing http request")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	if strings.Split(bodyString, ";")[0] != "redirecting to cas" {
		t.Errorf("Received an unexpected response from casauth")
	}

}
