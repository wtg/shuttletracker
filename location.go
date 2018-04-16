package shuttletracker

import (
	"time"
)

// Location represents information about a vehicle's location.
type Location struct {
	ID        int       `json:"id"`
	TrackerID string    `json:"tracker_id"`
	VehicleID int       `json:"vehicle_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Heading   float64   `json:"heading"`
	Speed     float64   `json:"speed"`
	Time      time.Time `json:"time"`
	RouteID   int       `json:"route_id"`
}

// LocationService is an interface for interacting with information about vehicle positions.
type LocationService interface {
	CreateLocation(location *Location) error
	DeleteLocationsBefore(before time.Time) (int, error)
	LocationsSince(vehicleID int, since time.Time) ([]*Location, error)
	LatestLocation(vehicleID int) (*Location, error)
}
