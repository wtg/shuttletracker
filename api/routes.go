package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/wtg/shuttletracker/log"

	"github.com/wtg/shuttletracker/model"
)

// RoutesHandler finds all of the routes in the database
func (api *API) RoutesHandler(w http.ResponseWriter, r *http.Request) {
	routes, err := api.db.GetRoutes()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		coord := model.Coord{
			Lat: c["lat"],
			Lng: c["lng"],
		}
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

	// Here do the interpolation
	// now we get the Segment for each segment ( this should be stored in database, just store it inside route for god sake)
	// Type conversions
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

	id := chi.URLParam(r, "id")
	log.Debugf("deleting", id)
	err := api.db.DeleteRoute(id)

	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Use this for importing a schedule
type sched struct {
	Times []model.Time `json:"times"`
	ID    string       `json:"id"`
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
	//log.Debugf("Able to decode stop.")

	// TODO: account for multiple routes
	route := model.Route{}
	routes := []model.Route{}
	// if len(stop.RouteIDS) == 1{
	// 	route, err := api.db.GetRoute(stop.RouteID)
	// } else {
		for i := 0; i < len(stop.RouteIDS); i++ {
			//routeString :=
			//log.Debugf(string(stop.RouteIDS[i]))
			route, err := api.db.GetRoute(string(stop.RouteIDS[i]))
			//log.Debugf("Route: " + route.Name)
			if err != nil {
				log.WithError(err).Error("Unable to get route.")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			routes = append(routes, route)
		}

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
	// We have to know the order of the stop and store a velocity vector into duration for the prediction
	// for i := range stop.RouteIDS {
	// 	routes[i].StopsID = append(routes[i].StopsID, stop.ID) // THIS REQUIRES the front end to have correct order << to be improved
	// 	err = api.db.ModifyRoute(&routes[i])
	// 	if err != nil {
	// 		log.WithError(err).Error("Unable to modify route.")
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	WriteJSON(w, stop)
}

func (api *API) StopsDeleteHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	log.Debugf("deleting", id)
	err := api.db.DeleteStop(id)

	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
