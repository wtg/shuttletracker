package migration

import "time"

type Coord struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

// represent a node on the graph
type RouteNode struct {
	ID         string    `json:"id"`
	Coordinate Coord     `json:"coord"`
	Created    time.Time `json:"created"`
	RouteID    string    `json:"routeid"`
	RouteOrder int       `json:"routeorder"` // this increments with 100 for better resolution of interpolation
	StopID     string    `json:"stopid"`
}

// represent a edge on the graph
type RouteEdge struct {
	ID       string    `json:"id"`
	Distance float64   `json:"distance"`
	Velocity float64   `json:"velocity"`
	StartID  string    `json:"start"`
	EndID    string    `json:"end"`
	Time     time.Time `json:"time"`
}

type ShuttleNode struct {
	ID         string `json:"id"`
	ShuttleID  string `json:"shuttleid"`
	Coordinate Coord  `json:"coord"`
}

// Vehicle represents an object being tracked.
type ShuttleInfo struct {
	ID      string    `json:"id"   bson:"vehicleID,omitempty"`
	Name    string    `json:"name" bson:"vehicleName"`
	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

// Route represent the information attached to a route
type RouteInfo struct {
	ID          string    `json:"id"             bson:"id"`
	Name        string    `json:"name"           bson:"name"`
	Description string    `json:"description"    bson:"description"`
	StartTime   string    `json:"startTime"      bson:"startTime"`
	EndTime     string    `json:"endTime"        bson:"endTime"`
	Enabled     bool      `json:"enabled,string" bson:"enabled"`
	Color       string    `json:"color"          bson:"color"`
	Width       int       `json:"width,string"   bson:"width"`
	Created     time.Time `json:"created"        bson:"created"`
	Updated     time.Time `json:"updated"        bson:"updated"`
}

// Stop indicates where a tracked object is scheduled to arrive
type StopInfo struct {
	ID          string `json:"id"             bson:"id"`
	Name        string `json:"name"           bson:"name"`
	Description string `json:"description"    bson:"description"`
	Address     string `json:"address"        bson:"address"`
	StartTime   string `json:"startTime"      bson:"startTime"`
	EndTime     string `json:"endTime"        bson:"endTime"`
	Enabled     bool   `json:"enabled,string" bson:"enabled"`
}
