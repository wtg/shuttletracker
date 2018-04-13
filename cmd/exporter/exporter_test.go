package exporter

import (
	"fmt"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/model"
	"testing"
)

// TestDump tests reading a database dump
func TestDump(t *testing.T) {
	d := database.Mock{}
	routes := []model.Route{
		model.Route{
			Name: "test",
		},
	}
  stops := []model.Stop{
    model.Stop{
      Name: "test",
    },
  }
  vehicles := []model.Vehicle{
    model.Vehicle{
      VehicleName: "test",
    },
  }
	// a = append(a,r)
	d.On("GetRoutes").Return(routes, nil)
  d.On("GetStops").Return(stops, nil)
  d.On("GetVehicles").Return(vehicles, nil)
  d.On("GetUsers").Return([]model.User{}, nil)
  d.On("GetMessages").Return([]model.AdminMessage{},nil)

  e := &Exporter{
    db: &d,
  }

  expected := &Dump{
    Routes: routes,
    Stops: stops,
    Vehicles:vehicles,
    Users: []model.User{},
    Messages: []model.AdminMessage{},
  }
  out,_ := e.read()
  s := fmt.Sprintf("%v", expected)
  s2 := fmt.Sprintf("%v", out)
	if(s != s2){
    t.Errorf("Expected dump did not match actual dump")
  }
  d.AssertExpectations(t)
}
