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
	Route     string    `json:"RouteID"     bson:"routeID"`
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

// Route is a temporary type alias during refactoring
type Route = shuttletracker.Route

// Time is a temporary type alias during refactoring
type Time = sttime.Time

// Coord is a temporary type alias during refactoring
type Coord = shuttletracker.Coord

// AdminMessage represents a message popup for the user from the site administrator
type AdminMessage struct {
	ID      int       `json:"id" 									bson:"id"`
	Type    string    `json:type								bson:"type"`
	Message string    `json:message						bson:"message"`
	Display bool      `json:display						bson:"display"`
	Created time.Time `json:created						bson:"created"`
}

// Stop indicates where a tracked object is scheduled to arrive
type Stop struct {
	ID           string  `json:"id"             bson:"id"`
	Name         string  `json:"name"           bson:"name"`
	Description  string  `json:"description"    bson:"description"`
	Lat          float64 `json:"lat,string"     bson:"lat"`
	Lng          float64 `json:"lng,string"     bson:"lng"`
	Address      string  `json:"address"        bson:"address"`
	StartTime    string  `json:"startTime"      bson:"startTime"`
	EndTime      string  `json:"endTime"        bson:"endTime"`
	Enabled      bool    `json:"enabled,string" bson:"enabled"`
	RouteID      string  `json:"routeId"        bson:"routeId"`
	SegmentIndex int     `json:"segmentindex"   bson:"segmentindex"`
}

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
