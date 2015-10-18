package tracking

import (
	"encoding/json"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Route represents a set of coordinates to draw a path on our tracking map
type Route struct {
	Name        string    `json:"name"           bson:"name"`
	Description string    `json:"description"    bson:"description"`
	StartTime   string    `json:"startTime"      bson:"startTime"`
	EndTime     string    `json:"endTime"        bson:"endTime"`
	Enabled     bool      `json:"enabled,string" bson:"enabled"`
	Color       string    `json:"color"          bson:"color"`
	Width       int       `json:"width,string"   bson:"width"`
	Coords      string    `json:"coords"         bson:"coords"`
	Created     time.Time `json:"created"        bson:"created"`
	Updated     time.Time `json:"updated"        bson:"updated"`
}

// Stop indicates where a tracked object is scheduled to arrive
type Stop struct {
	Name        string  `json:"name"        bson:"name"`
	Phonetic    string  `json:"phonetic"    bson:"phonetic"`
	Description string  `json:"description" bson:"description"`
	Address     string  `json:"address"     bson:"address"`
	TimeServed  string  `json:"timeServed"  bson:"timeServed"`
	Lat         float64 `json:"lat"         bson:"lat"`
	Lng         float64 `json:"lng"         bson:"lng"`
	Enabled     bool    `json:"enabled"     bson:"enabled"`
}

// RoutesHandler finds all of the routes in the database
func (App *App) RoutesHandler(w http.ResponseWriter, r *http.Request) {
	// Find all routes in database
	var routes []Route
	err := App.Routes.Find(bson.M{}).All(&routes)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each route to client as JSON
	WriteJSON(w, routes)
}

// RoutesCreateHandler adds a new route to the database
func (App *App) RoutesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new route object using request fields
	route := Route{}
	err := json.NewDecoder(r.Body).Decode(&route)

	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Store new route under routes collection
	err = App.Routes.Insert(&route)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
