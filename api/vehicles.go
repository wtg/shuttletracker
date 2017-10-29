package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/cas.v1"
	mgo "gopkg.in/mgo.v2"

	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/model"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var (
	lastUpdate time.Time
)

//DegreeToRadian converts degrees to radians.
func DegreeToRadian(val float64) float64 {
	ret := val * math.Pi / 180
	return ret
}

//ComputeDistance uses The Haversine Formula to return a distance between two coords in meters
func ComputeDistance(c1 model.Coord, c2 model.Coord) float64 {
	rad := 6371.0
	dLat := DegreeToRadian(c1.Lat - c2.Lat)
	dLng := DegreeToRadian(c1.Lat - c2.Lat)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(DegreeToRadian(c1.Lat))*math.Cos(DegreeToRadian(c2.Lat))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := rad * c * 1000 //distance in meters
	return d
}

func (api *API) GetCurrentStopForVehicle(vehicle model.Vehicle) (stop *model.Stop) {
	var update model.VehicleUpdate
	err := api.db.Updates.Find(bson.M{"vehicleID": vehicle.vehicleID}).Sort("-created").Limit(1).One(&update)
	if err != nil {
		log.WithError(err).Error("Unable to get vehicle update.")
		return nil
	}
	var stops []model.Stop
	err = api.db.Stops.Find(bson.M{}).All(&stops)

	var potentialStops []model.Stop

	for _, stop := range stops {
		clat, _ := strconv.ParseFloat(update.Lat,64)
		clng, _ := strconv.ParseFloat(update.Lng,64)
		stopCoord := model.Coord{
			Lat: stop.Lat,
			Lng: stop.Lng,
		}
		updateCoord := model.Coord{
			Lat: clat,
			Lng: clng,
		}
		fmt.Printf("%v\n",ComputeDistance(stopCoord,updateCoord))
		if(ComputeDistance(stopCoord,updateCoord) < 100){ // Closer than 100 meters
			potentialStops = append(potentialStops, stop);
		}
	}
	fmt.Printf("%v\n",potentialStops)

	return nil
}

// VehiclesHandler finds all the vehicles in the database.
func (api *API) VehiclesHandler(w http.ResponseWriter, r *http.Request) {
	// Find all vehicles in database
	var vehicles []model.Vehicle
	err := api.db.Vehicles.Find(bson.M{}).All(&vehicles)
	api.GetCurrentStopForVehicle(vehicles[0])
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
	err = api.db.Vehicles.Insert(&vehicle)
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

	err = api.db.Vehicles.Find(bson.M{"vehicleID": vehicle.VehicleID}).Sort("-created").Limit(1).One(&vehicle)
	vehicle.VehicleName = name
	vehicle.Enabled = enabled

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vehicle.Updated = time.Now()
	err = api.db.Vehicles.Update(bson.M{"vehicleID": vehicle.VehicleID}, vehicle)
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
	err := api.db.Vehicles.Remove(bson.M{"vehicleID": vars["id"]})
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Here's my view, keep every name the same meaning, otherwise, choose another.
// UpdatesHandler get the most recent update for each vehicle in the vehicles collection.
func (api *API) UpdatesHandler(w http.ResponseWriter, r *http.Request) {
	var vehicles []model.Vehicle

	// Query active and enabled vehicles
	err := api.db.Vehicles.Find(bson.M{"enabled": true}).All(&vehicles)
	if err != nil {
		log.WithError(err).Error("Unable to get enabled vehicles.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// allocate slice of capacity len(vehicles) and size zero
	updates := make([]model.VehicleUpdate, 0, len(vehicles))

	// Find most recent update in the last five minutes for each vehicle
	for _, vehicle := range vehicles {
		// here, huge waste of computational power, you record every shit inside the Updates table and using sort, I don't know what the hell is going on
		update := model.VehicleUpdate{}
		since := time.Now().Add(time.Minute * -5)
		err = api.db.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID, "created": bson.M{"$gt": since}}).Sort("-created").One(&update)
		if err == mgo.ErrNotFound {
			continue
		} else if err != nil {
			log.WithError(err).Error("Unable to get vehicle update.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		updates = append(updates, update)
	}

	// Convert updates to JSON
	WriteJSON(w, updates) // it's good to take some REST in our server :)
}

// UpdateMessageHandler generates a message about an update for a vehicle
func (App *API) UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	// For each vehicle/update, store message as a string
	var messages []string
	var message string
	var vehicles []model.Vehicle
	var update model.VehicleUpdate

	// Query all Vehicles
	err := App.db.Vehicles.Find(bson.M{}).All(&vehicles)
	// Handle errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Find recent updates and generate message
	for _, vehicle := range vehicles {
		// find 10 most recent records
		err := App.db.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID}).Sort("-created").Limit(1).One(&update)
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
