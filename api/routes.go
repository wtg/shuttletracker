package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/model"
	sttime "github.com/wtg/shuttletracker/time"
)

// RoutesHandler finds all of the routes in the database
func (api *API) RoutesHandler(w http.ResponseWriter, r *http.Request) {
	routes, err := api.ms.Routes()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	WriteJSON(w, routes)
}

// StopsHandler finds all of the route stops in the database
func (api *API) StopsHandler(w http.ResponseWriter, r *http.Request) {
	stops, err := api.ms.Stops()
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

	route := &shuttletracker.Route{
		Name:         routeData["name"],
		Description:  routeData["description"],
		TimeInterval: timeIntervals,
		Enabled:      enabled,
		Color:        routeData["color"],
		Width:        width,
		Coords:       coords,
		Created:      time.Now(),
		Updated:      time.Now()}
	err = api.ms.CreateRoute(route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RoutesDeleteHandler deletes a route from database
func (api *API) RoutesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.ms.DeleteRoute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Use this for importing a schedule
type sched struct {
	Times []model.Time `json:"times"`
	ID    int          `json:"id"`
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
	sttime.Sort(times.Times)
	route := &shuttletracker.Route{}
	route, err = api.ms.Route(times.ID)
	route.TimeInterval = times.Times

	err = api.ms.ModifyRoute(route)
	if err != nil {
		log.WithError(err).Error("Unable to store route into db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RoutesEditHandler Only handles editing enabled flag for now
func (api *API) RoutesEditHandler(w http.ResponseWriter, r *http.Request) {
	route := &shuttletracker.Route{}
	err := json.NewDecoder(r.Body).Decode(route)
	if err != nil {
		log.WithError(err).Error("Unable to decode route")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	en := route.Enabled
	route, err = api.ms.Route(route.ID)
	route.Enabled = en
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.ms.ModifyRoute(route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// StopsCreateHandler adds a new route stop to the database
func (api *API) StopsCreateHandler(w http.ResponseWriter, r *http.Request) {
	stop := &shuttletracker.Stop{}
	err := json.NewDecoder(r.Body).Decode(stop)
	if err != nil {
		log.WithError(err).Error("unable to decode stop")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.ms.CreateStop(stop)
	if err != nil {
		log.WithError(err).Error("unable to create stop")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, stop)
}

func (api *API) StopsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.ms.DeleteStop(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
