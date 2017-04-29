package tracking

import (
	"net/http"
	"gopkg.in/mgo.v2/bson"
//	log "github.com/Sirupsen/logrus"
	"strconv"
	"time"
	"math/big"
)

type LatestPosition struct {
	Longitude string `json:"longitude"`
	Latitude string `json:"latitude"`
	Timestamp time.Time `json:"timestamp"`
	Speed float64 `json:"speed"`
	Heading int `json:"heading"`
	Cardinal string `json:"cardinal_point"`
	StatusMessage *string `json:"public_status_message"` // this is a pointer so it defaults to null
}

type LegacyVehicle struct {
	Name string `json:"name"`
	ID int `json:"id"`
	LatestPosition LatestPosition `json:"latest_position"`
	Icon map[string]int `json:"icon"`
}

type LegacyVehicleContainer struct {
	Vehicle LegacyVehicle `json:"vehicle"`
}

type LegacyCoordinate struct {
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type LegacyRoute struct {
	Name string `json:"name"`
	Width int `json:"width"`
	ID big.Int `json:"id"`
	Color string `json:"color"`
	Coordinates []LegacyCoordinate `json:"coords"`
}

type LegacyStop struct {

}

type LegacyRoutesAndStopsContainer struct {
	Routes []LegacyRoute `json:"routes"`
	Stops []LegacyStop `json:"stops"`
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

		// legacy app expects vehicle ID to be a number...
		vehicleID, err := strconv.Atoi(vehicle.VehicleID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		latestPosition := LatestPosition{
			Longitude: update.Lng,
			Latitude: update.Lat,
			Heading: int(heading),
			Cardinal: cardinal,
			Speed: speed,
			Timestamp: update.Created,
		}

		legacy_vehicle := LegacyVehicle{
			Name: vehicle.VehicleName,
			ID: vehicleID,
			LatestPosition: latestPosition,
			Icon: map[string]int{"id": 1},
		}



		legacy_vehicles = append(legacy_vehicles, LegacyVehicleContainer{Vehicle: legacy_vehicle})
	}
	// Convert updates to JSON
	WriteJSON(w, legacy_vehicles)
}

func (App *App) LegacyRoutesHandler(w http.ResponseWriter, r *http.Request) {
	// Find all routes in database
	var routes []Route
	err := App.Routes.Find(bson.M{}).All(&routes)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert routes to legacy routes
	var legacyRoutes []LegacyRoute
	for _, route := range routes {
		// legacy app expects route ID to be a number, so we convert Mongo's base 16 ID to base 10 int
		var routeID big.Int
		routeID.SetString(route.ID, 16)

		// convert coordinates to legacy coordinates
		var coordinates []LegacyCoordinate
		for _, coordinate := range route.Coords {
			// convert from float to string 
			latitude := strconv.FormatFloat(coordinate.Lat, 'f', 5, 64)
			longitude := strconv.FormatFloat(coordinate.Lng, 'f', 5, 64)

			legacyCoordinate := LegacyCoordinate{
				Latitude: latitude,
				Longitude: longitude,
			}

			coordinates = append(coordinates, legacyCoordinate)
		}

		legacyRoute := LegacyRoute{
			Name: route.Name,
			Width: route.Width,
			Color: route.Color,
			ID: routeID,
			Coordinates: coordinates,
		}

		legacyRoutes = append(legacyRoutes, legacyRoute)
	}

	// Send to client as JSON
	routesAndStops := LegacyRoutesAndStopsContainer{
		Routes: legacyRoutes,
		Stops: nil,
	}
	WriteJSON(w, routesAndStops)
/*
	// Find all stops in databases
	var stops []Stop
	err := App.Stops.Find(bson.M{}).All(&stops)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each stop to client as JSON
	WriteJSON(w, stops)
	*/
}
