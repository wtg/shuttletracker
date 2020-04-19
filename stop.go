package shuttletracker

import (
	"errors"
	"time"
)

// Stop is a place where vehicles frequently stop.
type Stop struct {
	ID        int64     `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`

	// Name and Description are pointers because they may be nil.
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// StopService is an interface for interacting with Stops.
type StopService interface {
	Stop(id int64) (*Stop, error)
	Stops() ([]*Stop, error)
	CreateStop(stop *Stop) error
	DeleteStop(id int64) error
}

// ErrStopNotFound indicates that a Stop is not in the service.
var ErrStopNotFound = errors.New("Stop not found")
