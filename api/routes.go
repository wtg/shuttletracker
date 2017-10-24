package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	// MySQL driver
	"gopkg.in/cas.v1"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/wtg/shuttletracker/model"
	"gopkg.in/mgo.v2/bson"
)

//GetAllRoutes returns all of the routes in the routes collection
func (App *API) GetAllRoutes() (routes []model.Route) {
	err := App.db.Routes.Find(bson.M{}).All(&routes)
	if err != nil {
		log.Error(err)
	}
	return
}

//GetAllStops returns all of the stops in the stops collection
func (App *API) GetAllStops() (stops []model.Stop) {
	err := App.db.Stops.Find(bson.M{}).All(&stops)
	if err != nil {
		log.Error(err)
	}
	return
}

//GetStopById finds a stop in the database by it's id
func (App *API) GetStopById(id string) (stop model.Stop){

	err := App.db.Stops.Find(bson.M{"id": id}).One(&stop)
	if err != nil{
		log.Error(err)
	}
	return stop

}

//GetStopByName finds a stop by its name
func (App *API) GetStopById(name string) (stop model.Stop){

	err := App.db.Stops.Find(bson.M{"name": name}).One(&stop)
	if err != nil{
		log.Error(err)
	}
	return stop

}

// RoutesHandler sends a json of all of the routes in the database to the user
func (App *API) RoutesHandler(w http.ResponseWriter, r *http.Request) {
	routes := App.GetAllRoutes()

	WriteJSON(w, routes)
}

// StopsHandler finds all of the route stops in the database
func (App *API) StopsHandler(w http.ResponseWriter, r *http.Request) {
	stops := App.GetAllStops()

	WriteJSON(w, stops)
}


// RoutesCreateHandler adds a new route to the database
func (App *API) RoutesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new route object using request fields
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}
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
	coords := []model.Coord{}
	for _, c := range coordsData {
		coord := model.Coord{c["lat"], c["lng"]}
		coords = append(coords, coord)
	}

	// Here do the interpolation
	// now we get the Segment for each segment ( this should be stored in database, just store it inside route for god sake)
	fmt.Printf("Size of coordinates = %d", len(coords))
	// Type conversions
	enabled, _ := strconv.ParseBool(routeData["enabled"])
	width, _ := strconv.Atoi(routeData["width"])
	currentTime := time.Now()
	// Create a new route
	route := model.Route{
		ID:          bson.NewObjectId().Hex(),
		Name:        routeData["name"],
		Description: routeData["description"],
		StartTime:   routeData["startTime"],
		EndTime:     routeData["endTime"],
		Enabled:     enabled,
		Color:       routeData["color"],
		Width:       width,
		Coords:      coords,
		Created:     currentTime,
		Updated:     currentTime}
	// Store new route under routes collection
	err = App.db.Routes.Insert(&route)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//Deletes route from database
func (App *API) RoutesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}
	vars := mux.Vars(r)
	fmt.Printf(vars["id"])
	log.Debugf("deleting", vars["id"])
	err := App.db.Routes.Remove(bson.M{"id": vars["id"]})
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//RoutesEditHandler Only handles editing enabled flag for now
func (App *API) RoutesEditHandler(w http.ResponseWriter, r *http.Request) {
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}
	route := model.Route{}

	err := json.NewDecoder(r.Body).Decode(&route)
	en := route.Enabled
	if err != nil {
		fmt.Printf("lelel: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = App.db.Routes.Find(bson.M{"id": route.ID}).Sort("-created").Limit(1).One(&route)
	route.Enabled = en
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = App.db.Routes.Update(bson.M{"id": route.ID}, route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// StopsCreateHandler adds a new route stop to the database
func (App *API) StopsCreateHandler(w http.ResponseWriter, r *http.Request) {
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}
	// Create a new stop object using request fields
	stop := model.Stop{}
	err := json.NewDecoder(r.Body).Decode(&stop)
	stop.ID = bson.NewObjectId().Hex()
	route := model.Route{}
	err1 := App.db.Routes.Find(bson.M{"id": stop.RouteID}).One(&route)
	// Error handling

	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// We have to know the order of the stop and store a velocity vector into duration for the prediction
	route.StopsID = append(route.StopsID, stop.ID) // THIS REQUIRES the front end to have correct order << to be improved
	fmt.Println(route.StopsID)

	// Store new stop under stops collection
	err = App.db.Stops.Insert(&stop)
	// Error handling
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query := bson.M{"id": stop.RouteID}
	change := bson.M{"$set": bson.M{"availableroute": stop.SegmentIndex + 1, "stopsid": route.StopsID}}

	err = App.db.Routes.Update(query, change)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Closest Segment ID = " + strconv.Itoa(stop.SegmentIndex))
	WriteJSON(w, stop)
}

func (App *API) StopsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	vars := mux.Vars(r)
	log.Debugf("deleting", vars["id"])
	fmt.Printf(vars["id"])
	err := App.db.Stops.Remove(bson.M{"id": vars["id"]})
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
