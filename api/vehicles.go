package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"gopkg.in/cas.v1"

	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/model"
)

var (
	lastUpdate time.Time
)

// VehiclesHandler finds all the vehicles in the database.
func (api *API) VehiclesHandler(w http.ResponseWriter, r *http.Request) {
	// Find all vehicles in database
	vehicles, err := api.db.GetVehicles()

	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send each vehicle to client as JSON
	WriteJSON(w, vehicles)
}

// VehiclesCreateHandler adds a new vehicle to the database.
func (api *API) VehiclesCreateHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	// Create new vehicle object using request fields
	vehicle := model.Vehicle{}
	vehicle.Created = time.Now()
	vehicle.Updated = vehicle.Created
	vehicleData := json.NewDecoder(r.Body)
	err := vehicleData.Decode(&vehicle)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Store new vehicle under vehicles collection
	err = api.db.CreateVehicle(&vehicle)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (api *API) VehiclesEditHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	vehicle := model.Vehicle{}
	err := json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := vehicle.VehicleName
	enabled := vehicle.Enabled

	vehicle, err = api.db.GetVehicle(vehicle.VehicleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vehicle.VehicleName = name
	vehicle.Enabled = enabled
	vehicle.Updated = time.Now()

	err = api.db.ModifyVehicle(&vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) VehiclesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}
	// Delete vehicle from Vehicles collection
	id := chi.URLParam(r, "id")
	log.Debugf("deleting", id)
	err := api.db.DeleteVehicle(id)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Here's my view, keep every name the same meaning, otherwise, choose another.
// UpdatesHandler get the most recent update for each vehicle in the vehicles collection.
func (api *API) UpdatesHandler(w http.ResponseWriter, r *http.Request) {
	vehicles, err := api.db.GetEnabledVehicles()
	if err != nil {
		log.WithError(err).Error("Unable to get enabled vehicles.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// slice of capacity len(vehicles) and size zero
	updates := make([]model.VehicleUpdate, 0, len(vehicles))
	for _, vehicle := range vehicles {
		since := time.Now().Add(time.Minute * -5)
		vehicleUpdates, err := api.db.GetUpdatesForVehicleSince(vehicle.VehicleID, since)
		if err != nil {
			log.WithError(err).Error("Unable to get last vehicle update.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// if there is an update since the time, append it to all updates
		if len(vehicleUpdates) > 0 {
			updates = append(updates, vehicleUpdates[0])
		}
	}

	// Convert updates to JSON
	WriteJSON(w, updates) // it's good to take some REST in our server :)
}
