package api

import (
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

	vs.AssertExpectations(t)
	vs.AssertNumberOfCalls(t, "Vehicles", 1)
}
