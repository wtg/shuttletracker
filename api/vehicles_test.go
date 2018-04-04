package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/mock"
)

func TestVehiclesHandlerNoVehicles(t *testing.T) {
	vs := &mock.VehicleService{}
	vs.On("Vehicles").Return([]*shuttletracker.Vehicle{}, nil)

	api := API{
		vs: vs,
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("unable to create HTTP request: %s", err)
		return
	}

	api.VehiclesHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("got status code %d, expected 200", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("got Content-Type \"%s\", expected \"application/json\"", resp.Header.Get("Content-Type"))
	}

	var returnedVehicles []*shuttletracker.Vehicle
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&returnedVehicles)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if len(returnedVehicles) != 0 {
		t.Errorf("got %d vehicles, expected 0", len(returnedVehicles))
	}

	vs.AssertExpectations(t)
	vs.AssertNumberOfCalls(t, "Vehicles", 1)
}

func TestVehiclesHandlerTwoVehicles(t *testing.T) {
	vs := &mock.VehicleService{}
	vehicles := []*shuttletracker.Vehicle{
		{
			Name:    "Vehicle 1",
			Enabled: true,
		},
		{
			Name:    "Vehicle 2",
			Enabled: true,
		},
	}
	vs.On("Vehicles").Return(vehicles, nil)

	api := API{
		vs: vs,
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("unable to create HTTP request: %s", err)
		return
	}

	api.VehiclesHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("got status code %d, expected 200", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("got Content-Type \"%s\", expected \"application/json\"", resp.Header.Get("Content-Type"))
	}

	var returnedVehicles []*shuttletracker.Vehicle
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&returnedVehicles)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if len(returnedVehicles) != 2 {
		t.Errorf("got %d vehicles, expected 2", len(returnedVehicles))
	}

	for i := range vehicles {
		if *vehicles[i] != *returnedVehicles[i] {
			t.Errorf("got different vehicles at index %d: %+v expected %+v", i, returnedVehicles[i], vehicles[i])
		}
	}

	vs.AssertExpectations(t)
	vs.AssertNumberOfCalls(t, "Vehicles", 1)
}
