// Package model provides structs used by multiple packages within shuttletracker.
package model

import (
	"math/big"
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
	Segment   string    `json:"segment"     bson:"segment"` // the segment that a vehicle resides on
	Route			string		`json:"RouteID"			bson:"rotueID"`		
}

// Vehicle represents an object being tracked.
type Vehicle struct {
	VehicleID   string    `json:"vehicleID"   bson:"vehicleID,omitempty"`
	VehicleName string    `json:"vehicleName" bson:"vehicleName"`
	Created     time.Time `bson:"created"`
	Updated     time.Time `bson:"updated"`
	Active      bool      `json:"active"`
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

type LegacyVehicle struct {
	Name           string         `json:"name"`
	ID             int            `json:"id"`
	LatestPosition LatestPosition `json:"latest_position"`
	Icon           map[string]int `json:"icon"`
}

type LegacyVehicleContainer struct {
	Vehicle LegacyVehicle `json:"vehicle"`
}

type LegacyCoordinate struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type LegacyRoute struct {
	Name        string             `json:"name"`
	Width       int                `json:"width"`
	ID          big.Int            `json:"id"`
	Color       string             `json:"color"`
	Coordinates []LegacyCoordinate `json:"coords"`
}

type LegacyStopRoute struct {
	Name string  `json:"name"`
	ID   big.Int `json:"id"`
}

type LegacyStop struct {
	Name      string            `json:"name"`
	Longitude string            `json:"longitude"`
	Latitude  string            `json:"latitude"`
	ShortName string            `json:"short_name"`
	Routes    []LegacyStopRoute `json:"routes"`
}

type LegacyRoutesAndStopsContainer struct {
	Routes []LegacyRoute `json:"routes"`
	Stops  []LegacyStop  `json:"stops"`
}

// Coord represents a single lat/lng point used to draw routes
type Coord struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

// Route represents a set of coordinates to draw a path on our tracking map
type Route struct {
	ID             string    `json:"id"             bson:"id"`
	Name           string    `json:"name"           bson:"name"`
	Description    string    `json:"description"    bson:"description"`
	StartTime      string    `json:"startTime"      bson:"startTime"`
	EndTime        string    `json:"endTime"        bson:"endTime"`
	Enabled        bool      `json:"enabled,bool"	  bson:"enabled"`
	Color          string    `json:"color"          bson:"color"`
	Width          int       `json:"width,string"   bson:"width"`
	Coords         []Coord   `json:"coords"         bson:"coords"`
	Duration       []Segment `json:"duration"       bson:"duration"`
	StopsID        []string  `json:"stopsid"        bson:"stopsid"`
	AvailableRoute int       `json:"availableroute" bson:"availableroute"`
	Created        time.Time `json:"created"        bson:"created"`
	Updated        time.Time `json:"updated"        bson:"updated"`
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
