package migration

import (
	tracking "shuttle_tracking_2/tracking"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// represent the current location of a shuttle
type ShuttleCoordinate struct {
	ID string `json:"id"`
}

func Snap(R *mgo.Collection, S *mgo.Collection, C *mgo.Collection) {
	// conver the Route points into segment information
	var routes []tracking.Route
	err := R.Find(bson.M{}).All(routes)
	if err != nil {
		return
	}

	// insert the stop location into Route
	var stops []tracking.Stop
	err = S.Find(bson.M{}).All(&stops)
	if err != nil {
		return
	}
	for _, stop := range stops {
		route := RouteNode{
			ID:         string(bson.NewObjectId()),
			Coordinate: Coord{Lat: stop.Lat, Lng: stop.Lng},
			Created:    time.Now(),
			RouteID:    stop.RouteID,
			RouteOrder: 0,
			StopID:     stop.ID,
		}
		C.Insert(&route)
	}

}

// V will use information in Coord to generate Speed using Google API
func V(Coord *mgo.Collection, Velocity *mgo.Collection) {

}

// calculate the distance between any forward distance between two coordinates
func Distance() {

}

// calculate the time distance between any two points between two coordinates
func TimeDistance() {

}

// given a list of coordinates and a time T and a future time T', get a set of Coordinates that the shuttle could potentially arrive at.
func TravelDistance() {

}
