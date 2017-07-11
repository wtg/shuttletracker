package api

import (
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/model"
	"gopkg.in/mgo.v2/bson"
	"math/big"
	"net/http"
	"strconv"
)

func (App *API) LegacyVehiclesHandler(w http.ResponseWriter, r *http.Request) {
	// Query all Vehicles
	var vehicles []model.Vehicle
	err := App.db.Vehicles.Find(bson.M{}).All(&vehicles)
	// Handle errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Find recent updates for each vehicle
	var legacy_vehicles []model.LegacyVehicleContainer
	for _, vehicle := range vehicles {
		var update model.VehicleUpdate
		// here, huge waste of computational power, you record every shit inside the Updates table and using sort, I don't know what the hell is going on
		err := App.db.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID}).Sort("-created").Limit(1).One(&update)

		if err != nil {
			log.WithError(err).Error("Could not fetch vehicle updates.")
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

		// determine cardinal direction
		cardinal := CardinalDirection(&update.Heading)

		// legacy app expects vehicle ID to be a number...
		vehicleID, err := strconv.Atoi(vehicle.VehicleID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		latestPosition := model.LatestPosition{
			Longitude: update.Lng,
			Latitude:  update.Lat,
			Heading:   int(heading),
			Cardinal:  cardinal,
			Speed:     speed,
			Timestamp: update.Created,
		}

		legacy_vehicle := model.LegacyVehicle{
			Name:           vehicle.VehicleName,
			ID:             vehicleID,
			LatestPosition: latestPosition,
			Icon:           map[string]int{"id": 1},
		}

		legacy_vehicles = append(legacy_vehicles, model.LegacyVehicleContainer{Vehicle: legacy_vehicle})
	}
	// Convert updates to JSON
	WriteJSON(w, legacy_vehicles)
}

func (App *API) LegacyRoutesHandler(w http.ResponseWriter, r *http.Request) {
	// Find all routes in database
	var routes []model.Route
	err := App.db.Routes.Find(bson.M{}).All(&routes)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert routes to legacy routes
	var legacyRoutes []model.LegacyRoute
	for _, route := range routes {
		// legacy app expects route ID to be a number, so we convert Mongo's base 16 ID to base 10 int
		var routeID big.Int
		routeID.SetString(route.ID, 16)

		// convert coordinates to legacy coordinates
		var coordinates []model.LegacyCoordinate
		for _, coordinate := range route.Coords {
			// convert from float to string
			latitude := strconv.FormatFloat(coordinate.Lat, 'f', 5, 64)
			longitude := strconv.FormatFloat(coordinate.Lng, 'f', 5, 64)

			legacyCoordinate := model.LegacyCoordinate{
				Latitude:  latitude,
				Longitude: longitude,
			}

			coordinates = append(coordinates, legacyCoordinate)
		}

		legacyRoute := model.LegacyRoute{
			Name:        route.Name,
			Width:       route.Width,
			Color:       route.Color,
			ID:          routeID,
			Coordinates: coordinates,
		}

		legacyRoutes = append(legacyRoutes, legacyRoute)
	}

	// Find all stops in databases
	var stops []model.Stop
	err = App.db.Stops.Find(bson.M{}).All(&stops)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// convert stops to legacy stops
	var legacyStops []model.LegacyStop
	for _, stop := range stops {
		// see if this stop has already been created. this should probably use a map for faster lookup, but the data is small.
		found := false
		for i := range legacyStops {
			ls := &legacyStops[i]
			if ls.Name == stop.Name {
				// already created, so just append this route to the stop's routes instead of creating a duplicate

				// get route name
				var route model.Route
				err := App.db.Routes.Find(bson.M{"id": stop.RouteID}).One(&route)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// convert route ID to int
				var routeID big.Int
				routeID.SetString(route.ID, 16)

				legacyStopRoute := model.LegacyStopRoute{Name: route.Name, ID: routeID}
				ls.Routes = append(ls.Routes, legacyStopRoute)

				found = true
				break
			}
		}
		if found {
			continue
		}

		// convert from float to string
		latitude := strconv.FormatFloat(stop.Lat, 'f', 5, 64)
		longitude := strconv.FormatFloat(stop.Lng, 'f', 5, 64)

		// get route name
		var route model.Route
		err := App.db.Routes.Find(bson.M{"id": stop.RouteID}).One(&route)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// convert route ID to int
		var routeID big.Int
		routeID.SetString(route.ID, 16)

		legacyStopRoute := model.LegacyStopRoute{Name: route.Name, ID: routeID}
		legacyStopRoutes := []model.LegacyStopRoute{legacyStopRoute}

		legacyStop := model.LegacyStop{
			Name:      stop.Name,
			Longitude: longitude,
			Latitude:  latitude,
			ShortName: stop.Name,
			Routes:    legacyStopRoutes,
		}

		legacyStops = append(legacyStops, legacyStop)
	}

	// Send to client as JSON
	routesAndStops := model.LegacyRoutesAndStopsContainer{
		Routes: legacyRoutes,
		Stops:  legacyStops,
	}
	WriteJSON(w, routesAndStops)
}
