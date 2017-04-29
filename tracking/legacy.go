package tracking

import (
	"net/http"
	"gopkg.in/mgo.v2/bson"
	// log "github.com/Sirupsen/logrus"
	"strconv"
)

type LatestPosition struct {
	Longitude string `json:"longitude"`
	Latitude string `json:"latitude"`
	Timestamp string `json:"timestamp"`
	Speed float64 `json:"speed"`
	Heading int `json:"heading"`
	Cardinal string `json:"cardinal_point"`
	StatusMessage *string `json:"public_status_message"` // this is a pointer so it defaults to null
}

type LegacyVehicle struct {
	Name string `json:"name"`
	ID string `json:"id"`
	LatestPosition LatestPosition `json:"latest_position"`
	Icon map[string]int `json:"icon"`
}

type LegacyVehicleContainer struct {
	Vehicle LegacyVehicle `json:"vehicle"`
}

func (App *App) LegacyVehiclesHandler(w http.ResponseWriter, r *http.Request) {
	// Query all Vehicles
	var vehicles []Vehicle
	err := App.Vehicles.Find(bson.M{}).All(&vehicles)
	// Handle errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Find recent updates for each vehicle
	var legacy_vehicles []LegacyVehicleContainer
	for _, vehicle := range vehicles {
		var update VehicleUpdate
		// here, huge waste of computational power, you record every shit inside the Updates table and using sort, I don't know what the hell is going on
		err := App.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID}).Sort("-created").Limit(1).One(&update)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// turn an Update into a LegacyVehicle

		// convert speed from string (why????) to float as legacy API provided
		speed, err := strconv.ParseFloat(update.Speed, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// convert heading from string (why????) to float (and eventually int as legacy API provided)
		heading, err := strconv.ParseFloat(update.Heading, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// calculate cardinal direction
		var cardinal string
		if heading >= 315 || heading < 45 {
			cardinal = "North"
		} else if heading >= 45 && heading < 135 {
			cardinal = "East"
		} else if heading >= 135 && heading < 225 {
			cardinal = "South"
		} else if heading >= 225 && heading < 315 {
			cardinal = "West"
		}

		latestPosition := LatestPosition{
			Longitude: update.Lng,
			Latitude: update.Lat,
			Heading: int(heading),
			Cardinal: cardinal,
			Speed: speed,
		}

		legacy_vehicle := LegacyVehicle{
			Name: vehicle.VehicleName,
			ID: vehicle.VehicleID,
			LatestPosition: latestPosition,
			Icon: map[string]int{"id": 1},
		}



		legacy_vehicles = append(legacy_vehicles, LegacyVehicleContainer{Vehicle: legacy_vehicle})
	}
	// Convert updates to JSON
	WriteJSON(w, legacy_vehicles)
}
