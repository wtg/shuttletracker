package shuttletracker

import (
	"errors"
	"time"
)

// Route represents a set of coordinates to draw a path on our tracking map
type Route struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Color       string    `json:"color"`
	Width       int64     `json:"width"`
	StopIDs     []int64   `json:"stop_ids"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Points      []Point   `json:"points"`
}

// Point represents a latitude/longitude pair.
type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// RouteService is an interface for interacting with Routes.
type RouteService interface {
	Route(id int64) (*Route, error)
	Routes() ([]*Route, error)
	CreateRoute(route *Route) error
	DeleteRoute(id int64) error
	ModifyRoute(route *Route) error
}

// ErrRouteNotFound indicates that a Route is not in the service.
var ErrRouteNotFound = errors.New("Route not found")
