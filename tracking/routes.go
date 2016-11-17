package tracking

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"fmt"
	"io/ioutil"

	"gopkg.in/mgo.v2/bson"
)

// Coord represents a single lat/lng point used to draw routes
type Coord struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

// Route represents a set of coordinates to draw a path on our tracking map
type Route struct {
	ID          string    `json:"id"             bson:"_id,omitempty"`
	Name        string    `json:"name"           bson:"name"`
	Description string    `json:"description"    bson:"description"`
	StartTime   string    `json:"startTime"      bson:"startTime"`
	EndTime     string    `json:"endTime"        bson:"endTime"`
	Enabled     bool      `json:"enabled,string" bson:"enabled"`
	Color       string    `json:"color"          bson:"color"`
	Width       int       `json:"width,string"   bson:"width"`
	Coords      []Coord   `json:"coords"         bson:"coords"`
	Created     time.Time `json:"created"        bson:"created"`
	Updated     time.Time `json:"updated"        bson:"updated"`
}

// Stop indicates where a tracked object is scheduled to arrive
type Stop struct {
	ID          string `json:"id"             bson:"id"`
	Name        string `json:"name"           bson:"name"`
	Description string `json:"description"    bson:"description"`
	// position on map
	Lat     float64 `json:"lat,string"     bson:"lat"`
	Lng     float64 `json:"lng,string"     bson:"lng"`
	Address string  `json:"address"        bson:"address"`

	StartTime string `json:"startTime"      bson:"startTime"`
	EndTime   string `json:"endTime"        bson:"endTime"`
	Enabled   bool   `json:"enabled,string" bson:"enabled"`
	RouteID   string `json:"routeId"        bson:"routeId"`
}

type MapPoint struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
type MapResponsePoint struct {
	Location      MapPoint `json:"location"`
	OriginalIndex int      `json:"originalIndex,omitempty"`
	PlaceID       string   `json:"placeId"`
}
type MapResponse struct {
	SnappedPoints []MapResponsePoint
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

// StopsHandler finds all of the route stops in the database
func (App *App) StopsHandler(w http.ResponseWriter, r *http.Request) {
	// Find all stops in database
	var stops []Stop
	err := App.Stops.Find(bson.M{}).All(&stops)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each stop to client as JSON
	WriteJSON(w, stops)
}

func Interpolate(coords []Coord, key string) []Coord {
	// make request
	prefix := "https://roads.googleapis.com/v1/snapToRoads?"
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString("path=")
	for i, coord := range coords {
		buffer.WriteString(strconv.FormatFloat(coord.Lat, 'f', 10, 64))
		buffer.WriteString(",")
		buffer.WriteString(strconv.FormatFloat(coord.Lng, 'f', 10, 64))
		if i < len(coords)-1 {
			buffer.WriteString("|")
		}
	}
	buffer.WriteString("&interpolate=true&key=")
	buffer.WriteString(key)
	fmt.Print("sending request: " + buffer.String())
	resp, err := http.Get(buffer.String())
	if err != nil {
		fmt.Errorf("Error Not valid response from Google API")
		return nil
	}
	defer resp.Body.Close()
	fmt.Print(resp.Header)
	body, err := ioutil.ReadAll(resp.Body)

	mapResponse := MapResponse{}
	json.Unmarshal(body, &mapResponse)
	result := []Coord{}
	for _, location := range mapResponse.SnappedPoints {
		currentLocation := Coord{
			Lat: float64(location.Location.Latitude),
			Lng: float64(location.Location.Longitude),
		}
		result = append(result, currentLocation)
	}
	return result
}

// RoutesCreateHandler adds a new route to the database
func (App *App) RoutesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new route object using request fields
	var routeData map[string]string
	var coordsData []map[string]float64
	// Decode route details
	err := json.NewDecoder(r.Body).Decode(&routeData)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Unmarshal route coordinates
	err = json.Unmarshal([]byte(routeData["coords"]), &coordsData)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Create a Coord from each set of input coordinates
	coords := []Coord{}
	for _, c := range coordsData {
		coord := Coord{c["lat"], c["lng"]}
		coords = append(coords, coord)
	}
	// Here do the interpolation
	coords = Interpolate(coords, App.Config.GoogleMapAPIKey)
	fmt.Printf("Size of coordinates = %d", len(coords))
	// Type conversions
	enabled, _ := strconv.ParseBool(routeData["enabled"])
	width, _ := strconv.Atoi(routeData["width"])
	currentTime := time.Now()
	// Create a new route
	route := Route{
		string(bson.NewObjectId()),
		routeData["name"],
		routeData["description"],
		routeData["startTime"],
		routeData["endTime"],
		enabled,
		routeData["color"],
		width,
		coords,
		currentTime,
		currentTime}
	// Store new route under routes collection
	err = App.Routes.Insert(&route)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Deletes route from database
func (App *App) RoutesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debugf("deleting", vars["id"])
	err := App.Routes.Remove(bson.M{"_id": bson.ObjectIdHex(vars["id"])})
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// StopsCreateHandler adds a new route stop to the database
func (App *App) StopsCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new stop object using request fields
	stop := Stop{}
	err := json.NewDecoder(r.Body).Decode(&stop)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Store new stop under stops collection
	err = App.Stops.Insert(&stop)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (App *App) StopsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debugf("deleting", vars["id"])
	err := App.Stops.Remove(bson.M{"name": vars["id"]})
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
