package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/wtg/shuttletracker"
)

type VehicleService struct {
	mock.Mock
}

// CreateVehicle creates a vehicle.
func (vs *VehicleService) CreateVehicle(vehicle *shuttletracker.Vehicle) error {
	args := vs.Called(vehicle)
	return args.Error(0)
}

// DeleteVehicle deletes a vehicle.
func (vs *VehicleService) DeleteVehicle(vehicleID int) error {
	args := vs.Called(vehicleID)
	return args.Error(0)
}

// GetVehicle gets a vehicle.
func (vs *VehicleService) Vehicle(vehicleID int) (*shuttletracker.Vehicle, error) {
	args := vs.Called(vehicleID)
	return args.Get(0).(*shuttletracker.Vehicle), args.Error(1)
}

func (vs *VehicleService) VehicleWithTrackerID(trackerID int) (*shuttletracker.Vehicle, error) {
	args := vs.Called(trackerID)
	return args.Get(0).(*shuttletracker.Vehicle), args.Error(1)
}

// GetVehicles gets all vehicles.
func (vs *VehicleService) Vehicles() ([]*shuttletracker.Vehicle, error) {
	args := vs.Called()
	return args.Get(0).([]*shuttletracker.Vehicle), args.Error(1)
}

// GetEnabledVehicles gets all enabled vehicles.
func (vs *VehicleService) EnabledVehicles() ([]*shuttletracker.Vehicle, error) {
	args := vs.Called()
	return args.Get(0).([]*shuttletracker.Vehicle), args.Error(1)
}

// ModifyVehicle modifies a vehicle.
func (vs *VehicleService) ModifyVehicle(vehicle *shuttletracker.Vehicle) error {
	args := vs.Called(vehicle)
	return args.Error(0)
}
