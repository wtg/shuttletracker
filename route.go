package shuttletracker

import (
	"time"

	sttime "github.com/wtg/shuttletracker/time"
)

// Route represents a set of coordinates to draw a path on our tracking map
type Route struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	TimeInterval []sttime.Time `json:"intervals"`
	Enabled      bool          `json:"enabled"`
	Active       bool          `json:"active"`
	Color        string        `json:"color"`
	Width        int           `json:"width"`
	Coords       []Coord       `json:"coords"`
	StopsID      []string      `json:"stopsid"`
	Created      time.Time     `json:"created"        bson:"created"`
	Updated      time.Time     `json:"updated"        bson:"updated"`
}

// Coord represents a single lat/lng point used to draw routes
type Coord struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

// RouteService is an interface for interacting with Routes.
type RouteService interface {
	Route(id int) (*Route, error)
	Routes() ([]*Route, error)
	CreateRoute(route *Route) error
	DeleteRoute(id int) error
	ModifyRoute(route *Route) error
}
