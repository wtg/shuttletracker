package shuttletracker

import (
	"time"
)

// Position represents information about a vehicle's location.
type Position struct {
	ID int `json:"vehicleID"   bson:"vehicleID,omitempty"`

	VehicleID string    `json:"vehicleID"   bson:"vehicleID,omitempty"`
	Latitude  float64   `json:"lat"         bson:"lat"`
	Longitude float64   `json:"lng"         bson:"lng"`
	Heading   int64     `json:"heading"     bson:"heading"`
	Speed     float64   `json:"speed"       bson:"speed"`
	Time      time.Time `json:"time"        bson:"time"`
}

// LocationService is an interface for interacting with information about vehicle positions.
type LocationService interface {
	LatestPosition(vehicleID int) (*Position, error)
}
