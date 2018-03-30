package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/model"
)

func TestCardinalDirection(t *testing.T) {
	table := [][]string{
		{"0", "North"},
		{"45", "North-East"},
		{"90", "East"},
		{"135", "South-East"},
		{"180", "South"},
		{"225", "South-West"},
		{"270", "West"},
		{"315", "North-West"},
		{"this isn't a direction lol", "North"},
	}

	for _, testCase := range table {
		direction := CardinalDirection(&testCase[0])
		expected := testCase[1]
		if direction != expected {
			t.Errorf("Got %v, expected %v.", direction, expected)
		}
	}
}

func TestVehiclesHandler(t *testing.T) {
	type testCase struct {
		method string
		path   string
	}
	cases := []testCase{
		{
			method: "GET",
			path:   "/vehicles",
		},
	}

	for _, c := range cases {
		cfg := Config{}
		db := &database.Mock{}
		db.On("GetVehicles").Return([]model.Vehicle{}, nil)

		api, err := New(cfg, db)
		if err != nil {
			t.Errorf("got error '%s', expected nil", err)
			return
		}

		server := httptest.NewServer(api.handler)
		defer server.Close()
		client := http.Client{}

		url := server.URL + c.path
		req, err := http.NewRequest(c.method, url, nil)
		if err != nil {
			t.Errorf("unable to create HTTP request: %s", err)
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("unable to do request: %s", err)
			continue
		}

		if resp.StatusCode != 200 {
			t.Logf("%+v", req)
			t.Logf("%+v", resp)
			t.Errorf("%s %s returned status code %d, expected 200", c.method, url, resp.StatusCode)
		}
	}
}
