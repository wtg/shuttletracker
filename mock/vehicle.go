package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/wtg/shuttletracker"
)

// VehicleService implements a mock of shuttletracker.VehicleService.
type VehicleService struct {
	mock.Mock
}

// CreateVehicle creates a Vehicle.
func (vs *VehicleService) CreateVehicle(vehicle *shuttletracker.Vehicle) error {
	args := vs.Called(vehicle)
	return args.Error(0)
}

// DeleteVehicle deletes a Vehicle.
func (vs *VehicleService) DeleteVehicle(vehicleID int64) error {
	args := vs.Called(vehicleID)
	return args.Error(0)
}

// Vehicle gets a Vehicle.
func (vs *VehicleService) Vehicle(vehicleID int64) (*shuttletracker.Vehicle, error) {
	args := vs.Called(vehicleID)
	return args.Get(0).(*shuttletracker.Vehicle), args.Error(1)
}

// VehicleWithTrackerID returns a Vehicle with the specified tracker ID.
func (vs *VehicleService) VehicleWithTrackerID(trackerID string) (*shuttletracker.Vehicle, error) {
	args := vs.Called(trackerID)
	return args.Get(0).(*shuttletracker.Vehicle), args.Error(1)
}

// Vehicles gets all Vehicles.
func (vs *VehicleService) Vehicles() ([]*shuttletracker.Vehicle, error) {
	args := vs.Called()
	return args.Get(0).([]*shuttletracker.Vehicle), args.Error(1)
}

// EnabledVehicles gets all enabled Vehicles.
func (vs *VehicleService) EnabledVehicles() ([]*shuttletracker.Vehicle, error) {
	args := vs.Called()
	return args.Get(0).([]*shuttletracker.Vehicle), args.Error(1)
}

// ModifyVehicle modifies a Vehicle.
func (vs *VehicleService) ModifyVehicle(vehicle *shuttletracker.Vehicle) error {
	args := vs.Called(vehicle)
	return args.Error(0)
}
