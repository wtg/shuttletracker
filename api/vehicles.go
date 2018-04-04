package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/cas.v1"

	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/model"
)

var (
	lastUpdate time.Time
)

// VehiclesHandler returns all the vehicles.
func (api *API) VehiclesHandler(w http.ResponseWriter, r *http.Request) {
	vehicles, err := api.vs.Vehicles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, vehicles)
}

// VehiclesCreateHandler adds a new vehicle.
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
	err = api.vs.CreateVehicle(&vehicle)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (api *API) VehiclesEditHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	vehicle := &model.Vehicle{}
	err := json.NewDecoder(r.Body).Decode(vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := vehicle.Name
	enabled := vehicle.Enabled
	trackerID := vehicle.TrackerID

	vehicle, err = api.vs.Vehicle(vehicle.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vehicle.Name = name
	vehicle.Enabled = enabled
	vehicle.TrackerID = trackerID
	vehicle.Updated = time.Now()

	err = api.vs.ModifyVehicle(vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) VehiclesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.vs.DeleteVehicle(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdatesHandler gets the most recent update for each enabled vehicle.
func (api *API) UpdatesHandler(w http.ResponseWriter, r *http.Request) {
	vehicles, err := api.vs.EnabledVehicles()
	if err != nil {
		log.WithError(err).Error("Unable to get enabled vehicles.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// slice of capacity len(vehicles) and size zero
	updates := make([]model.VehicleUpdate, 0, len(vehicles))
	for _, vehicle := range vehicles {
		since := time.Now().Add(time.Minute * -5)
		vehicleUpdates, err := api.db.GetUpdatesForVehicleSince(vehicle.ID, since)
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
