package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/cas.v1"

	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/model"

	"github.com/gorilla/mux"
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
	vars := mux.Vars(r)
	log.Debugf("deleting", vars["id"])
	err := api.db.DeleteVehicle(vars["id"])
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

// UpdateMessageHandler generates a message about an update for a vehicle
func (api *API) UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	// For each vehicle/update, store message as a string
	var messages []string
	var message string

	// Query all Vehicles
	vehicles, err := api.db.GetVehicles()
	// Handle errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Find recent updates and generate message
	for _, vehicle := range vehicles {
		// find 10 most recent records
		update, err := api.db.GetLastUpdateForVehicle(vehicle.VehicleID)
		if err == nil {
			// Use first 4 char substring of update.Speed
			speed := update.Speed
			if len(speed) > 4 {
				speed = speed[0:4]
			}

			// Convert last updated time to local timezone
			loc, err := time.LoadLocation("America/New_York")
			if err != nil {
				log.WithError(err).Error("Could not load time zone information.")
				continue
			}
			lastUpdate := update.Created.In(loc).Format("3:04:05pm")

			message = fmt.Sprintf("<b>%s</b><br/>Traveling %s at<br/> %s mph as of %s", vehicle.VehicleName, CardinalDirection(&update.Heading), speed, lastUpdate)
			messages = append(messages, message)
		}
	}
	// Convert to JSON
	WriteJSON(w, messages)
}

// CardinalDirection returns the cardinal direction of a vehicle's heading.
func CardinalDirection(h *string) string {
	heading, err := strconv.ParseFloat(*h, 64)
	if err != nil {
		log.WithError(err).Error("Unable to parse float")
		return "North"
	}
	switch {
	case (heading >= 22.5 && heading < 67.5):
		return "North-East"
	case (heading >= 67.5 && heading < 112.5):
		return "East"
	case (heading >= 112.5 && heading < 157.5):
		return "South-East"
	case (heading >= 157.5 && heading < 202.5):
		return "South"
	case (heading >= 202.5 && heading < 247.5):
		return "South-West"
	case (heading >= 247.5 && heading < 292.5):
		return "West"
	case (heading >= 292.5 && heading < 337.5):
		return "North-West"
	default:
		return "North"
	}
}
