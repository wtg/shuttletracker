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
	StopIDs      []int64       `json:"stop_ids"`
	Created      time.Time     `json:"created"`
	Updated      time.Time     `json:"updated"`
	Points       []Point       `json:"points"`
}

// Point represents a latitude/longitude pair.
type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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

//SetRouteActiveStatus determines if a given route is active based on its schedule intervals and the time given, then updates the object in the parameter
func SetRouteActiveStatus(r *Route, t time.Time) {

	//This is a time offset, to ensure routes are activated on the minute they are assigned activate
	var currentTime sttime.Time
	currentTime.FromTime(t)
	currentTime.Day = time.Now().Weekday()
	state := -1

	if r.TimeInterval == nil || len(r.TimeInterval) == 1 {
		state = 1
	}
	for idx, val := range r.TimeInterval {
		//If it is the last in the time list (latest time for the week) use this index
		if idx >= len(r.TimeInterval)-1 {
			state = val.State
			break
		} else {
			if currentTime.After(val) && currentTime.After(r.TimeInterval[idx+1]) {
				continue
			}
			state = val.State
			break
		}
	}

	r.Active = (state == 1 || state == -1)

}
