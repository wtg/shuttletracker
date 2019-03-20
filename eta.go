package shuttletracker

import (
	"time"
)

// VehicleETA contains information about a Vehicle, its current Route, and its
// ETAs to Stops.
type VehicleETA struct {
	VehicleID int64     `json:"vehicle_id"`
	RouteID   int64     `json:"route_id"`
	StopETAs  []StopETA `json:"stop_etas"`
	Updated   time.Time `json:"updated"`
}

// StopETA represents a time when a Vehicle is expected to arrive at a Stop.
type StopETA struct {
	StopID   int64     `json:"stop_id"`
	ETA      time.Time `json:"eta"`
	Arriving bool      `json:"arriving"`
}

// ETAService is an interface for interacting with vehicle estimated times of arrival.
type ETAService interface {
	Subscribe(func(VehicleETA))
	CurrentETAs() map[int64]VehicleETA
}
