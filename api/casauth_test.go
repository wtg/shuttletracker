package api

import (
	"github.com/go-chi/chi"
	gc "gopkg.in/cas.v2"

	"github.com/wtg/shuttletracker/auth"
	"github.com/wtg/shuttletracker/database"

	"io/ioutil"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCasUnauthenticated(t *testing.T) {
	url, _ := url.Parse("https://cas.example.com/")
	c := gc.NewClient(&gc.Options{
		URL: url,
	})
	client := auth.CAS{
		CAS : c,
	}
	db := &database.Mock{}
	httpcli := http.Client{}
	cli := casClient{
		cas: &client,
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

func TestCasAuthenticated(t *testing.T) {

	client := &auth.Mock{}
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
	db.On("UserExists","lyonj4").Return(true,nil)

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
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if (bodyString != "test"){
		t.Errorf("Response did not come through, authenticaiton failure")
	}
	_ = req
	_ = httpcli
	_ = resp
	_ = err

}
func TestCasAuthenticatedBadUser(t *testing.T) {

	client := &auth.Mock{}
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
	db.On("UserExists","lyonj4").Return(false,nil)

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
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if (bodyString != "unauthenticated\n"){
		t.Errorf("Response should be unauthenticated")
	}
	_ = req
	_ = httpcli
	_ = resp
	_ = err

}
