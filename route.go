package shuttletracker

import (
	"time"

	sttime "github.com/wtg/shuttletracker/time"
)

// Route represents a set of coordinates to draw a path on our tracking map
type Route struct {
	ID             string        `json:"id"             bson:"id"`
	Name           string        `json:"name"           bson:"name"`
	Description    string        `json:"description"    bson:"description"`
	TimeInterval   []sttime.Time `json:"intervals"			 bson:"intervals"`
	Enabled        bool          `json:"enabled,bool"	 bson:"enabled"`
	Active         bool          `json:"active,bool"	 bson:"enabled"`
	Color          string        `json:"color"          bson:"color"`
	Width          int           `json:"width,string"   bson:"width"`
	Coords         []Coord       `json:"coords"         bson:"coords"`
	StopsID        []string      `json:"stopsid"        bson:"stopsid"`
	AvailableRoute int           `json:"availableroute" bson:"availableroute"`
	Created        time.Time     `json:"created"        bson:"created"`
	Updated        time.Time     `json:"updated"        bson:"updated"`
}

// Coord represents a single lat/lng point used to draw routes
type Coord struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

// RouteService is an interface for interacting with Routes.
type RouteService interface {
}
