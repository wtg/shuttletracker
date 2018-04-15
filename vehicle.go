package shuttletracker

import (
	"errors"
	"time"
)

// ErrVehicleNotFound indicates that a Vehicle is not in the service.
var ErrVehicleNotFound = errors.New("Vehicle not found")

// Vehicle represents an object being tracked.
type Vehicle struct {
	ID        int       `json:"id"   bson:"vehicleID,omitempty"`
	Name      string    `json:"name" bson:"vehicleName"`
	Created   time.Time `json:"created" bson:"created"`
	Updated   time.Time `json:"updated" bson:"updated"`
	Enabled   bool      `json:"enabled"     bson:"enabled"`
	TrackerID string    `json:"tracker_id"`
}

// VehicleService is an interface for interacting with Vehicles.
type VehicleService interface {
	Vehicle(id int) (*Vehicle, error)
	VehicleWithTrackerID(id string) (*Vehicle, error)
	Vehicles() ([]*Vehicle, error)
	EnabledVehicles() ([]*Vehicle, error)
	CreateVehicle(vehicle *Vehicle) error
	DeleteVehicle(id int) error
	ModifyVehicle(vehicle *Vehicle) error
}
