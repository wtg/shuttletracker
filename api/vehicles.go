package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
)

var (
	lastUpdate time.Time
)

// VehiclesHandler returns all the vehicles.
func (api *API) VehiclesHandler(w http.ResponseWriter, r *http.Request) {
	vehicles, err := api.ms.Vehicles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, vehicles)
}

// VehiclesCreateHandler adds a new vehicle.
func (api *API) VehiclesCreateHandler(w http.ResponseWriter, r *http.Request) {
	vehicle := shuttletracker.Vehicle{}
	vehicleData := json.NewDecoder(r.Body)
	err := vehicleData.Decode(&vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.ms.CreateVehicle(&vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (api *API) VehiclesEditHandler(w http.ResponseWriter, r *http.Request) {
	vehicle := &shuttletracker.Vehicle{}
	err := json.NewDecoder(r.Body).Decode(vehicle)
	if err != nil {
		log.WithError(err).Error("unable to decode vehicle")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := vehicle.Name
	enabled := vehicle.Enabled
	trackerID := vehicle.TrackerID
	vehicle, err = api.ms.Vehicle(vehicle.ID)
	if err != nil {
		log.WithError(err).Error("unable to retrieve vehicle")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vehicle.Name = name
	vehicle.Enabled = enabled
	vehicle.TrackerID = trackerID

	err = api.ms.ModifyVehicle(vehicle)
	if err != nil {
		log.WithError(err).Error("unable to modify vehicle")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) VehiclesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.ms.DeleteVehicle(id)
	if err != nil {
		if err == shuttletracker.ErrVehicleNotFound {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// UpdatesHandler gets the most recent update for each enabled vehicle.
func (api *API) UpdatesHandler(w http.ResponseWriter, r *http.Request) {
	vehicles, err := api.ms.EnabledVehicles()
	if err != nil {
		log.WithError(err).Error("Unable to get enabled vehicles.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// slice of capacity len(vehicles) and size zero
	updates := make([]*shuttletracker.Location, 0, len(vehicles))
	for _, vehicle := range vehicles {
		since := time.Now().Add(time.Minute * -5)
		vehicleUpdates, err := api.ms.LocationsSince(vehicle.ID, since)
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
