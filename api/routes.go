package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
	"sort"

	// MySQL driver
	"gopkg.in/cas.v1"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/wtg/shuttletracker/model"
)

//TODO: Move this to updater
//RouteIsActive determines if the current time means a route should be active or not
func (api *API) RouteIsActive(r *model.Route) (bool){
	currentTime := time.Now()
	day := currentTime.Weekday()
	p := fmt.Println
	state := -1
	for idx, val := range r.TimeInterval{
			if(day == val.Day){
				if (currentTime.After(val.Time) && !currentTime.After(r.TimeInterval[idx+1].Time)){
					state = val.State
					p("state",state)
				}
				}else if (day > val.Day){
					if(r.TimeInterval[idx+1].Day < day){
						state = val.State
						p("state",state)

						}else{
							if(currentTime.After(val.Time)&&r.TimeInterval[idx+1].Time.After(currentTime)){
								state = val.State
								p("state",state)

							}
						}
					}
				}


	return false
}

// RoutesHandler finds all of the routes in the database
func (api *API) RoutesHandler(w http.ResponseWriter, r *http.Request) {
	// Find all routes in database
	routes, err := api.db.GetRoutes()


	for idx,_ := range routes{

		timeIntervals := []model.WeekTime{}

		//Temp times for testing
		on := time.Now().AddDate(0,0,-20)
		off := time.Now().Add(-1*time.Minute)

		//Strip the date from the time by bringing the day up to today - Do this when we add a new range
		val := -(int(on.Sub(time.Now().Truncate(24*time.Hour)).Hours()/24)-1)
		on = on.AddDate(0,0,val)

		val = -(int(off.Sub(time.Now().Truncate(24*time.Hour)).Hours()/24)-1)
		off = off.AddDate(0,0,val)

		testTime := model.WeekTime{
			Day: time.Sunday,
			Time: on,
			State: 1,
		}


		testTime3 := model.WeekTime{
			Day: time.Monday,
			Time: on,
			State: 4,
		}

		testTime2 := model.WeekTime{
			Day: time.Monday,
			Time: off,
			State: 2,
		}

		timeIntervals = append(timeIntervals,testTime2)
		timeIntervals = append(timeIntervals,testTime3)
		timeIntervals = append(timeIntervals,testTime)
		sort.Sort(model.ByTime(timeIntervals))

		routes[idx].TimeInterval = timeIntervals

		api.RouteIsActive(&routes[idx])
	}
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each route to client as JSON
	WriteJSON(w, routes)
}

// StopsHandler finds all of the route stops in the database
func (api *API) StopsHandler(w http.ResponseWriter, r *http.Request) {
	// Find all stops in databases
	stops, err := api.db.GetStops()
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each stop to client as JSON
	WriteJSON(w, stops)
}

// compute distance between two coordinates and return a value
func ComputeDistance(c1 model.Coord, c2 model.Coord) float64 {
	return float64(math.Sqrt(math.Pow(c1.Lat-c2.Lat, 2) + math.Pow(c1.Lng-c2.Lng, 2)))
}

func ComputeDistanceMapPoint(c1 model.MapPoint, c2 model.MapPoint) float64 {
	return float64(math.Sqrt(math.Pow(c1.Latitude-c2.Latitude, 2) + math.Pow(c1.Longitude-c2.Longitude, 2)))
}

// RoutesCreateHandler adds a new route to the database
func (api *API) RoutesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new route object using request fields
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
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
	timeIntervals := []model.WeekTime{}

	//This will never have an error
  form := "3:04pm";
	midnight,_ := time.Parse(form,"12:00am")

	testTime := model.WeekTime{
		Day: time.Monday,
		Time: midnight,
		State: 1,
	}

	timeIntervals = append(timeIntervals,testTime)


	// Create a new route
	route := model.Route{
		Name:        routeData["name"],
		Description: routeData["description"],
		TimeInterval:     timeIntervals,
		Enabled:     enabled,
		Color:       routeData["color"],
		Width:       width,
		Coords:      coords,
		Created:     currentTime,
		Updated:     currentTime}
	// Store new route under routes collection
	err = api.db.CreateRoute(&route)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// RoutesDeleteHandler deletes a route from database
func (api *API) RoutesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}
	vars := mux.Vars(r)
	fmt.Printf(vars["id"])
	log.Debugf("deleting", vars["id"])
	err := api.db.DeleteRoute(vars["id"])
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RoutesEditHandler Only handles editing enabled flag for now
func (api *API) RoutesEditHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}
	route := model.Route{}

	err := json.NewDecoder(r.Body).Decode(&route)
	if err != nil {
		log.WithError(err).Error("Unable to decode route")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	en := route.Enabled

	route, err = api.db.GetRoute(route.ID)
	route.Enabled = en
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = api.db.ModifyRoute(&route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// StopsCreateHandler adds a new route stop to the database
func (api *API) StopsCreateHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	// Create a new stop object using request fields
	stop := model.Stop{}
	err := json.NewDecoder(r.Body).Decode(&stop)
	if err != nil {
		log.WithError(err).Error("Unable to decode stop.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	route, err := api.db.GetRoute(stop.RouteID)
	if err != nil {
		log.WithError(err).Error("Unable to get route.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Store new stop under stops collection
	err = api.db.CreateStop(&stop)
	// Error handling
	if err != nil {
		log.WithError(err).Error("Unable to create stop.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// We have to know the order of the stop and store a velocity vector into duration for the prediction
	route.StopsID = append(route.StopsID, stop.ID) // THIS REQUIRES the front end to have correct order << to be improved
	err = api.db.ModifyRoute(&route)
	if err != nil {
		log.WithError(err).Error("Unable to modify route.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, stop)
}

// StopsDeleteHandler deletes a Stop.
func (api *API) StopsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	vars := mux.Vars(r)
	log.Debugf("deleting", vars["id"])
	fmt.Printf(vars["id"])
	err := api.db.DeleteStop(vars["id"])
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
