package api

import (
	"github.com/go-chi/chi"

	"github.com/wtg/shuttletracker/auth"
	"github.com/wtg/shuttletracker/mock"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCasUnauthenticated(t *testing.T) {
	url, _ := url.Parse("https://cas.example.com/")

	us := &mock.UserService{}

	cli := CreateCASClient(url, us, true)
	httpcli := http.Client{}

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
	us := &mock.UserService{}
	httpcli := http.Client{}

	cli := InjectMocks(client, us, true)
	r := chi.NewRouter()
	r.Use(cli.casauth)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})
	us.On("UserExists", "lyonj4").Return(true, nil)

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

	if bodyString != "test" {
		t.Errorf("Response did not come through, authenticaiton failure")
	}
	us.AssertExpectations(t)
	_ = req
	_ = httpcli
	_ = resp
	_ = err

}
func TestCasAuthenticatedBadUser(t *testing.T) {

	client := &auth.Mock{}
	us := &mock.UserService{}
	httpcli := http.Client{}
	cli := InjectMocks(client, us, true)

	r := chi.NewRouter()
	r.Use(cli.casauth)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})
	us.On("UserExists", "lyonj4").Return(false, nil)

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
	if bodyString != "unauthenticated\n" {
		t.Errorf("Response should be unauthenticated")
	}
	us.AssertExpectations(t)

	_ = req
	_ = httpcli
	_ = resp
	_ = err

}
