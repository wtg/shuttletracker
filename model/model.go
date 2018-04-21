// Package model provides structs used by multiple packages within shuttletracker.
package model

import (
	"time"

	"github.com/wtg/shuttletracker"
	sttime "github.com/wtg/shuttletracker/time"
)

// Vehicle is a temporary type alias during refactoring
type Vehicle = shuttletracker.Vehicle

// VehicleUpdate represents a single position observed for a Vehicle.
type VehicleUpdate struct {
	VehicleID string    `json:"vehicleID"   bson:"vehicleID,omitempty"`
	Lat       string    `json:"lat"         bson:"lat"`
	Lng       string    `json:"lng"         bson:"lng"`
	Heading   string    `json:"heading"     bson:"heading"`
	Speed     string    `json:"speed"       bson:"speed"`
	Lock      string    `json:"lock"        bson:"lock"`
	Time      string    `json:"time"        bson:"time"`
	Date      string    `json:"date"        bson:"date"`
	Status    string    `json:"status"      bson:"status"`
	Created   time.Time `json:"created"     bson:"created"`
	Route     int       `json:"RouteID"     bson:"routeID"`
}

// Status contains a detailed message on the tracked object's status.
type Status struct {
	Public  bool      `bson:"public"`
	Message string    `json:"message" bson:"message"`
	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

type LatestPosition struct {
	Longitude     string    `json:"longitude"`
	Latitude      string    `json:"latitude"`
	Timestamp     time.Time `json:"timestamp"`
	Speed         float64   `json:"speed"`
	Heading       int       `json:"heading"`
	Cardinal      string    `json:"cardinal_point"`
	StatusMessage *string   `json:"public_status_message"` // this is a pointer so it defaults to null
}

// Time is a temporary type alias during refactoring
type Time = sttime.Time

// Coord is a temporary type alias during refactoring
type Coord = shuttletracker.Coord

type MapPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type MapResponsePoint struct {
	Location      MapPoint `json:"location"`
	OriginalIndex int      `json:"originalIndex,omitempty"`
	PlaceID       string   `json:"placeId"`
}
type MapResponse struct {
	SnappedPoints []MapResponsePoint
}
