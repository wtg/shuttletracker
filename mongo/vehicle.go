package mongo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/wtg/shuttletracker"
)

// VehicleService is a MongoDB implementation of shuttletracker's VehicleService.
type VehicleService struct {
	c mgo.Collection
}

// CreateVehicle creates a Vehicle.
func (v *VehicleService) CreateVehicle(vehicle *shuttletracker.Vehicle) error {
	return v.c.Insert(&vehicle)
}

// DeleteVehicle deletes a Vehicle by its ID.
func (v *VehicleService) DeleteVehicle(vehicleID string) error {
	return v.c.Remove(bson.M{"vehicleID": vehicleID})
}

// GetVehicle returns a Vehicle by its ID.
func (v *VehicleService) GetVehicle(vehicleID string) (*shuttletracker.Vehicle, error) {
	var vehicle *shuttletracker.Vehicle
	err := v.c.Find(bson.M{"vehicleID": vehicleID}).One(&vehicle)
	if err == mgo.ErrNotFound {
		return vehicle, shuttletracker.ErrVehicleNotFound
	}
	return vehicle, err
}

// GetVehicles returns all Vehicles.
func (v *VehicleService) GetVehicles() ([]*shuttletracker.Vehicle, error) {
	var vehicles []*shuttletracker.Vehicle
	err := v.c.Find(bson.M{}).All(&vehicles)
	return vehicles, err
}

// GetEnabledVehicles returns all Vehicles that are enabled.
func (v *VehicleService) GetEnabledVehicles() ([]*shuttletracker.Vehicle, error) {
	var vehicles []*shuttletracker.Vehicle
	err := v.c.Find(bson.M{"enabled": true}).All(&vehicles)
	return vehicles, err
}

// ModifyVehicle updates a Vehicle by its ID.
func (v *VehicleService) ModifyVehicle(vehicle *shuttletracker.Vehicle) error {
	return v.c.Update(bson.M{"vehicleID": vehicle.ID}, vehicle)
}
