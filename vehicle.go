package shuttletracker

import (
	"errors"
	"time"
)

// ErrVehicleNotFound indicates that a Vehicle is not in the service.
var ErrVehicleNotFound = errors.New("Vehicle not found")

// Vehicle represents an object being tracked.
type Vehicle struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Enabled   bool      `json:"enabled"`
	TrackerID string    `json:"tracker_id"`
}

// VehicleService is an interface for interacting with Vehicles.
type VehicleService interface {
	Vehicle(id int64) (*Vehicle, error)
	VehicleWithTrackerID(id string) (*Vehicle, error)
	Vehicles() ([]*Vehicle, error)
	EnabledVehicles() ([]*Vehicle, error)
	CreateVehicle(vehicle *Vehicle) error
	DeleteVehicle(id int64) error
	ModifyVehicle(vehicle *Vehicle) error
}
