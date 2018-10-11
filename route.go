package shuttletracker

import (
	"errors"
	"time"
)

// Route represents a set of coordinates to draw a path on our tracking map
type Route struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Enabled     bool          `json:"enabled"`
	Color       string        `json:"color"`
	Width       int64         `json:"width"`
	StopIDs     []int64       `json:"stop_ids"`
	Created     time.Time     `json:"created"`
	Updated     time.Time     `json:"updated"`
	Points      []Point       `json:"points"`
	Active      bool          `json:"active"`
	Schedule    RouteSchedule `json:"schedule"`
}

// RouteActiveInterval represents a time interval during which a Route is active.
type RouteActiveInterval struct {
	ID        int64        `json:"id"`
	RouteID   int64        `json:"route_id"`
	StartDay  time.Weekday `json:"start_day"`
	StartTime time.Time    `json:"start_time"`
	EndDay    time.Weekday `json:"end_day"`
	EndTime   time.Time    `json:"end_time"`
}

// RouteSchedule represents multiple time intervals during which a Route is active.
type RouteSchedule []RouteActiveInterval

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
