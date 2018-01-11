// Package model provides structs used by multiple packages within shuttletracker.
package model

import (
	"time"
)

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

// Vehicle represents an object being tracked.
type Vehicle struct {
	VehicleID   string    `json:"vehicleID"   bson:"vehicleID,omitempty"`
	VehicleName string    `json:"vehicleName" bson:"vehicleName"`
	Created     time.Time `bson:"created"`
	Updated     time.Time `bson:"updated"`
	Enabled     bool      `json:"enabled"     bson:"enabled"`
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

// Coord represents a single lat/lng point used to draw routes
type Coord struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

//WeekTime represents a time of the week that routes will be turned on or off, include a time and day
//The state variable is intentionally vague so that it can be used for several different applications
type WeekTime struct {
	Day   time.Weekday `json:"day"     bson:"day"`
	Time  time.Time    `json:"time"    bson:"time"`
	State int          `json:"on"   bson:"on"`
}

type ByTime []WeekTime

func (a ByTime) Len() int      { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool {
	if a[i].Day < a[j].Day {
		return true
	} else if a[i].Day > a[j].Day {
		return false
	} else {
		if a[j].Time.After(a[i].Time) {
			return true
		} else {
			return false
		}
	}
}

// Route represents a set of coordinates to draw a path on our tracking map
type Route struct {
	ID             string     `json:"id"             bson:"id"`
	Name           string     `json:"name"           bson:"name"`
	Description    string     `json:"description"    bson:"description"`
	TimeInterval   []WeekTime `json:"intervals"			 bson:"intervals"`
	Enabled        bool       `json:"enabled,bool"	 bson:"enabled"`
	Active         bool       `json:"active,bool"	 bson:"enabled"`
	Color          string     `json:"color"          bson:"color"`
	Width          int        `json:"width,string"   bson:"width"`
	Coords         []Coord    `json:"coords"         bson:"coords"`
	Duration       []Segment  `json:"duration"       bson:"duration"`
	StopsID        []string   `json:"stopsid"        bson:"stopsid"`
	AvailableRoute int        `json:"availableroute" bson:"availableroute"`
	Created        time.Time  `json:"created"        bson:"created"`
	Updated        time.Time  `json:"updated"        bson:"updated"`
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

type MapDistanceMatrixDuration struct {
	Value int    `json:"value"`
	Text  string `json:"text"`
}

type MapDistanceMatrixDistance struct {
	Value int    `json:"value"`
	Text  string `json:"text"`
}

type MapDistanceMatrixElement struct {
	Status   string                    `json:"status"`
	Duration MapDistanceMatrixDuration `json:"duration"`
	Distance MapDistanceMatrixDistance `json:"distance"`
}

type MapDistanceMatrixElements struct {
	Elements []MapDistanceMatrixElement `json:"elements"`
}
type MapDistanceMatrixResponse struct {
	Status               string                      `json:"status"`
	OriginAddresses      []string                    `json:"origin_addresses"`
	DestinationAddresses []string                    `json:"destination_addresses"`
	Rows                 []MapDistanceMatrixElements `json:"rows"`
}

type Segment struct {
	ID       string   `json:"id"`
	Start    MapPoint `json:"origin"`
	End      MapPoint `json:"destination"`
	Distance float64  `json:"distance"`
	Duration float64  `json:"duration"`
}
