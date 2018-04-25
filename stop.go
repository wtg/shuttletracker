package shuttletracker

import "time"

// Stop is a place where vehicles frequently stop.
type Stop struct {
	ID        int       `json:"id"`
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
	Stops() ([]*Stop, error)
	CreateStop(stop *Stop) error
	DeleteStop(id int) error
}
