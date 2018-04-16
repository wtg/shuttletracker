package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/mock"
)

func TestVehiclesHandlerNoVehicles(t *testing.T) {
	ms := &mock.ModelService{}
	ms.VehicleService.On("Vehicles").Return([]*shuttletracker.Vehicle{}, nil)

	api := API{
		ms: ms,
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

	ms.VehicleService.AssertExpectations(t)
	ms.VehicleService.AssertNumberOfCalls(t, "Vehicles", 1)
}

func TestVehiclesHandlerTwoVehicles(t *testing.T) {
	ms := &mock.ModelService{}
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
	ms.VehicleService.On("Vehicles").Return(vehicles, nil)

	api := API{
		ms: ms,
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

	ms.VehicleService.AssertExpectations(t)
	ms.VehicleService.AssertNumberOfCalls(t, "Vehicles", 1)
}

func TestVehiclesCreateHandler(t *testing.T) {
	ms := &mock.ModelService{}
	vehicle := &shuttletracker.Vehicle{
		Name:      "Vehicle 2",
		Enabled:   true,
		TrackerID: "2",
	}
	ms.VehicleService.On("CreateVehicle", vehicle).Return(nil)

	api := API{
		ms: ms,
	}

	body := &bytes.Buffer{}
	enc := json.NewEncoder(body)
	err := enc.Encode(vehicle)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Errorf("unable to create HTTP request: %s", err)
		return
	}

	w := httptest.NewRecorder()
	api.VehiclesCreateHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("got status code %d, expected 200", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if len(respBody) != 0 {
		t.Errorf("got body length %d, expected 0", len(respBody))
	}

	ms.VehicleService.AssertExpectations(t)
	ms.VehicleService.AssertNumberOfCalls(t, "CreateVehicle", 1)
}

func TestVehiclesEditHandler(t *testing.T) {
	ms := &mock.ModelService{}
	vehicleTime := time.Now()
	existingVehicle := &shuttletracker.Vehicle{
		ID:        4,
		Name:      "Vehicle 2",
		Enabled:   true,
		TrackerID: "2",
		Created:   vehicleTime,
	}
	changedVehicle := &shuttletracker.Vehicle{
		ID:        4,
		Name:      "Vehicle 2 changed",
		Enabled:   false,
		TrackerID: "3",
		Created:   vehicleTime,
	}
	ms.VehicleService.On("Vehicle", 4).Return(existingVehicle, nil)

	// Because the handler sets the Updated field to time.Now(), we have to accept everything
	// here and then check later that the method actually got the struct with fields we expected.
	ms.VehicleService.On("ModifyVehicle", "mock.Anything").Return(nil)

	api := API{
		ms: ms,
	}

	body := &bytes.Buffer{}
	enc := json.NewEncoder(body)
	err := enc.Encode(changedVehicle)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Errorf("unable to create HTTP request: %s", err)
		return
	}

	w := httptest.NewRecorder()
	api.VehiclesEditHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("got status code %d, expected 200", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if len(respBody) != 0 {
		t.Errorf("got body length %d, expected 0", len(respBody))
	}

	ms.VehicleService.AssertExpectations(t)
	ms.VehicleService.AssertNumberOfCalls(t, "Vehicle", 1)
	ms.VehicleService.AssertNumberOfCalls(t, "ModifyVehicle", 1)

	for _, call := range ms.VehicleService.Calls {
		if call.Method == "ModifyVehicle" {
			if len(call.Arguments) != 1 {
				t.Errorf("got %d arguments, expected 1", len(call.Arguments))
				continue
			}
			argVehicle, ok := call.Arguments[0].(*shuttletracker.Vehicle)
			if !ok {
				t.Error("expected ModifyVehicle to be called with *Vehicle argument")
			}

			if argVehicle.Created != changedVehicle.Created {
				t.Error("got unexpected vehicle.Created value")
			}
			if argVehicle.Name != changedVehicle.Name {
				t.Error("got unexpected vehicle.Name value")
			}
			if argVehicle.TrackerID != changedVehicle.TrackerID {
				t.Error("got unexpected vehicle.TrackerID value")
			}
			if argVehicle.ID != changedVehicle.ID {
				t.Error("got unexpected vehicle.ID value")
			}
			if argVehicle.Enabled != changedVehicle.Enabled {
				t.Error("got unexpected vehicle.Enabled value")
			}
			break
		}
	}
}

func TestVehiclesDeleteHandler(t *testing.T) {
	ms := &mock.ModelService{}
	vehicleID := 7
	ms.VehicleService.On("DeleteVehicle", vehicleID).Return(nil)

	api := API{
		ms: ms,
	}

	req, err := http.NewRequest("DELETE", "/?id="+strconv.Itoa(vehicleID), nil)
	if err != nil {
		t.Errorf("unable to create HTTP request: %s", err)
		return
	}

	w := httptest.NewRecorder()
	api.VehiclesDeleteHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("got status code %d, expected 200", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if len(respBody) != 0 {
		t.Errorf("got body length %d, expected 0", len(respBody))
	}

	ms.VehicleService.AssertExpectations(t)
	ms.VehicleService.AssertNumberOfCalls(t, "DeleteVehicle", 1)
}
