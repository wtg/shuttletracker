package shuttletracker

import (
	"errors"
	"time"
)

// ErrVehicleNotFound indicates that a Vehicle is not in the database.
var ErrVehicleNotFound = errors.New("Vehicle not found")

// Vehicle represents an object being tracked.
type Vehicle struct {
	ID        int       `json:"vehicleID"   bson:"vehicleID,omitempty"`
	Name      string    `json:"vehicleName" bson:"vehicleName"`
	Created   time.Time `bson:"created"`
	Updated   time.Time `bson:"updated"`
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
