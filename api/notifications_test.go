package api

// TODO: fix these tests

/*
import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wtg/shuttletracker/database"
)

func TestNew(t *testing.T) {
	type testCase struct {
		method string
		path   string
	}
	cases := []testCase{
		{
			method: "GET",
			path:   "/",
		},
		{
			method: "GET",
			path:   "/updates",
		},
	}

	cfg := Config{}
	db := &database.Mock{}

	api, err := New(cfg, db)
	if err != nil {
		t.Errorf("got error '%s', expected nil", err)
	}

	server := httptest.NewServer(api.handler)
	client := server.Client()
	for _, c := range cases {
		url := server.URL + c.path
		t.Log(url)
		req, err := http.NewRequest(c.method, url, nil)
		if err != nil {
			t.Errorf("unable to create HTTP request: %s", err)
			continue
		}
		_, err = client.Do(req)
		if err != nil {
			t.Errorf("unable to do request: %s", err)
			continue
		}
	}
}
*/
