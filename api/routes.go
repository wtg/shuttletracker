package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"
	// MySQL driver
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/wtg/shuttletracker/model"
)

//RouteIsActive determines if the current time means a route should be active or not
func (api *API) RouteIsActive(r *model.Route) bool {
	var currentTime model.Time
	currentTime.FromTime(time.Now())
	currentTime.Day = time.Now().Weekday()
	state := -1

	if r.TimeInterval == nil || len(r.TimeInterval) == 1 {
		state = 1
	}
	for idx, val := range r.TimeInterval {
		//If it is the last in the time list (latest time for the week) use this index
		if idx >= len(r.TimeInterval)-1 {
			state = val.State
			break
		} else {
			if currentTime.After(val) && currentTime.After(r.TimeInterval[idx+1]) {
				continue
			}
			state = val.State
			break
		}
	}
	route := model.Route{}
	if api.db != nil {
		r, err := api.db.GetRoute(r.ID)
		route = r
		if err != nil {
			return false
		}
	}
	//If we cannot determine a state for some reason default to active
	route.Active = (state == 1 || state == -1)
	if api.db != nil {
		err := api.db.ModifyRoute(&route)
		if err != nil {
			return false
		}
	}
	return (state == 1 || state == -1)

}

// RoutesHandler finds all of the routes in the database
func (api *API) RoutesHandler(w http.ResponseWriter, r *http.Request) {
	routes, err := api.db.GetRoutes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for idx := range routes {
		api.RouteIsActive(&routes[idx])
	}
	WriteJSON(w, routes)
}

// StopsHandler finds all of the route stops in the database
func (api *API) StopsHandler(w http.ResponseWriter, r *http.Request) {
	stops, err := api.db.GetStops()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	WriteJSON(w, stops)
}

func combineCoords(coordsData *[]map[string]float64) []model.Coord {
	coords := []model.Coord{}
	for _, c := range *coordsData {
		coord := model.Coord{c["lat"], c["lng"]}
		coords = append(coords, coord)
	}
	return coords
}

// RoutesCreateHandler adds a new route to the database
func (api *API) RoutesCreateHandler(w http.ResponseWriter, r *http.Request) {
	var routeData map[string]string
	var coordsData []map[string]float64
	err := json.NewDecoder(r.Body).Decode(&routeData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.Unmarshal([]byte(routeData["coords"]), &coordsData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	coords := combineCoords(&coordsData)

	enabled, _ := strconv.ParseBool(routeData["enabled"])
	width, _ := strconv.Atoi(routeData["width"])
	timeIntervals := []model.Time{}

	route := model.Route{
		Name:         routeData["name"],
		Description:  routeData["description"],
		TimeInterval: timeIntervals,
		Enabled:      enabled,
		Color:        routeData["color"],
		Width:        width,
		Coords:       coords,
		Created:      time.Now(),
		Updated:      time.Now()}
	err = api.db.CreateRoute(&route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RoutesDeleteHandler deletes a route from database
func (api *API) RoutesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debugf("[ROUTE DELETE:]", vars["id"])
	err := api.db.DeleteRoute(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Use this for importing a schedule
type sched struct {
	Times []model.Time `json:times`
	ID    string       `json:id`
}

// RoutesScheduler Allows for route active times to be set
func (api *API) RoutesScheduler(w http.ResponseWriter, r *http.Request) {
	var times sched
	err := json.NewDecoder(r.Body).Decode(&times)
	if err != nil {
		log.WithError(err).Error("Unable to decode route")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sort.Sort(model.ByTime(times.Times))
	route := model.Route{}
	route, err = api.db.GetRoute(times.ID)
	route.TimeInterval = times.Times
	err = api.db.ModifyRoute(&route)
	if err != nil {
		log.WithError(err).Error("Unable to store route into db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RoutesEditHandler Only handles editing enabled flag for now
func (api *API) RoutesEditHandler(w http.ResponseWriter, r *http.Request) {
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
	var stop model.Stop
	err := json.NewDecoder(r.Body).Decode(&stop)
	if err != nil {
		log.WithError(err).Error("Unable to decode stop")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	route, err := api.db.GetRoute(stop.RouteID)
	if err != nil {
		log.WithError(err).Error("Unable to get route.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.CreateStop(&stop)
	if err != nil {
		log.WithError(err).Error("Unable to create stop.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	route.StopsID = append(route.StopsID, stop.ID)
	err = api.db.ModifyRoute(&route)
	if err != nil {
		log.WithError(err).Error("Unable to modify route.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, stop)
}

func (api *API) StopsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debugf("deleting", vars["id"])
	err := api.db.DeleteStop(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
