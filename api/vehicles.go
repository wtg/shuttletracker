package api

import (
	"encoding/json"
	"fmt"
	"gopkg.in/cas.v1"
	"net/http"
	"strconv"
	"time"

	"github.com/wtg/shuttletracker/log"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var (
	lastUpdate time.Time
)

//  represents a single position observed for a Vehicle from the data feed.
type VehicleUpdate struct {
	VehicleID string    `json:"vehicleID"   bson:"vehicleID,omitempty"`
	Lat       string    `json:"lat"         bson:"lat"`
	Lng       string    `json:"lng"         bson:"lng"`
	Heading   string    `json:"heading"     bson:"heading"`
	Speed     string    `json:"speed"       bson:"speed"`
	Lock      string    `json:"lock"        bson:"lock"`
	Time      string    `json:"time"        bson:"time"`
	Date      string    `json:"date"        bson:"date"`
	Status    string    `json:"status"      bson:"status"`
	Created   time.Time `json:"created"     bson:"created"`
	Segment   string    `json:"segment"     bson:"segment"` // the segment that a vehicle resides on
}

// Vehicle represents an object being tracked.
type Vehicle struct {
	VehicleID   string    `json:"vehicleID"   bson:"vehicleID,omitempty"`
	VehicleName string    `json:"vehicleName" bson:"vehicleName"`
	Created     time.Time `bson:"created"`
	Updated     time.Time `bson:"updated"`
	Active bool `json:"active"`
}

// Status contains a detailed message on the tracked object's status
type Status struct {
	Public  bool      `bson:"public"`
	Message string    `json:"message" bson:"message"`
	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}


// VehiclesHandler finds all the vehicles in the database.
func (App *API) VehiclesHandler(w http.ResponseWriter, r *http.Request) {

	// Find all vehicles in database
	var vehicles []Vehicle
	err := App.db.Vehicles.Find(bson.M{}).All(&vehicles)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send each vehicle to client as JSON
	WriteJSON(w, vehicles)
}

// VehiclesCreateHandler adds a new vehicle to the database.
func (App *API) VehiclesCreateHandler(w http.ResponseWriter, r *http.Request) {
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	// Create new vehicle object using request fields
	vehicle := Vehicle{}
	vehicle.Created = time.Now()
	vehicle.Updated = time.Now()
	vehicleData := json.NewDecoder(r.Body)
	err := vehicleData.Decode(&vehicle)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Store new vehicle under vehicles collection
	err = App.db.Vehicles.Insert(&vehicle)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (App *API) VehiclesEditHandler(w http.ResponseWriter, r *http.Request) {

}

func (App *API) VehiclesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	// Delete vehicle from Vehicles collection
	vars := mux.Vars(r)
	log.Debugf("deleting", vars["id"])
	err := App.db.Vehicles.Remove(bson.M{"vehicleID": vars["id"]})
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Here's my view, keep every name the same meaning, otherwise, choose another.
// UpdatesHandler get the most recent update for each vehicle in the vehicles collection.
func (App *API) UpdatesHandler(w http.ResponseWriter, r *http.Request) {
	// Store updates for each vehicle
	var vehicles []Vehicle
	var updates []VehicleUpdate
	var update VehicleUpdate
	var vehicleUpdates []VehicleUpdate
	// Query all Vehicles
	err := App.db.Vehicles.Find(bson.M{}).All(&vehicles)
	// Handle errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Find recent updates for each vehicle
	for _, vehicle := range vehicles {
		// here, huge waste of computational power, you record every shit inside the Updates table and using sort, I don't know what the hell is going on
		err := App.db.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID}).Sort("-created").Limit(20).All(&vehicleUpdates);
		update = vehicleUpdates[0]

		if err == nil {
			count := 0.0
			speed := 0.0
			for i, elem := range vehicleUpdates{
				if(time.Since(elem.Created).Minutes() < 5){
					val,_ := strconv.ParseFloat(vehicleUpdates[i].Speed,64);
					speed += val
					count += 1;
				}
			}
			if(count > 0 && speed/count > 5){
				updates = append(updates, update)
			}
		}
	}

	// Convert updates to JSON
	WriteJSON(w, updates) // it's good to take some REST in our server :)
}

// UpdateMessageHandler generates a message about an update for a vehicle
func (App *API) UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	// For each vehicle/update, store message as a string
	var messages []string
	var message string
	var vehicles []Vehicle
	var update VehicleUpdate

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
			//nextArrival := GetArrivalTime(&update, App.Routes, App.Stops)
			message = fmt.Sprintf("<b>%s</b><br/>Traveling %s at<br/> %s mph as of %s", vehicle.VehicleName, CardinalDirection(&update.Heading), speed, lastUpdate.Format("3:04:05pm") /*, nextArrival*/)
			messages = append(messages, message)
		}
	}
	// Convert to JSON
	WriteJSON(w, messages)
}

// CardinalDirection figures out the cardinal direction of a vehicle's heading
func CardinalDirection(h *string) string {
	heading, err := strconv.ParseFloat(*h, 64)
	if err != nil {
		fmt.Println("ERROR", err.Error())
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

