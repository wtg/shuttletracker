package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
)

func (api *API) ETAHandler(w http.ResponseWriter, r *http.Request) {
	etas := api.etaManager.CurrentETAs()
	err := WriteJSON(w, etas)
	if err != nil {
		return
	}
}

// RoutesHandler finds all of the routes in the database
func (api *API) RoutesHandler(w http.ResponseWriter, r *http.Request) {
	routes, err := api.ms.Routes()
	if err != nil {
		log.WithError(err).Error("unable to get routes")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, routes)
}

// StopsHandler finds all of the route stops in the database
func (api *API) StopsHandler(w http.ResponseWriter, r *http.Request) {
	stops, err := api.ms.Stops()
	if err != nil {
		log.WithError(err).Error("unable to get stops")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, stops)
}

// RoutesCreateHandler adds a new route to the database
func (api *API) RoutesCreateHandler(w http.ResponseWriter, r *http.Request) {
	route := &shuttletracker.Route{}
	err := json.NewDecoder(r.Body).Decode(route)
	if err != nil {
		log.WithError(err).Error("unable to decode route")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = api.ms.CreateRoute(route)
	if err != nil {
		log.WithError(err).Error("unable to create route")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RoutesDeleteHandler deletes a route from database
func (api *API) RoutesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = api.ms.DeleteRoute(id)
	if err != nil {
		if err == shuttletracker.ErrRouteNotFound {
			http.Error(w, "Route not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	sched := route.Schedule
	route, err = api.ms.Route(route.ID)
	route.Enabled = en
	route.Schedule = sched
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
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.ms.DeleteStop(id)
	if err != nil {
		if err == shuttletracker.ErrStopNotFound {
			http.Error(w, "Stop not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// func (api *API) UnmarshalJSON(data []byte) error
