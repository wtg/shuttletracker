package tracking

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// VehicleUpdate represents a single position observed for a Vehicle from the data feed.
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
}

// Vehicle represents an object being tracked.
type Vehicle struct {
	VehicleID   string    `json:"vehicleID"   bson:"vehicleID,omitempty"`
	VehicleName string    `json:"vehicleName" bson:"vehicleName"`
	Created     time.Time `bson:"created"`
	Updated     time.Time `bson:"updated"`
}

// Status contains a detailed message on the tracked object's status
type Status struct {
	Public  bool      `bson:"public"`
	Message string    `json:"message" bson:"message"`
	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

var (
	// Match each API field with any number (+)
	//   of the previous expressions (\d digit, \. escaped period, - negative number)
	//   Specify named capturing groups to store each field from data feed
	dataRe    = regexp.MustCompile(`(?P<id>Vehicle ID:([\d\.]+)) (?P<lat>lat:([\d\.-]+)) (?P<lng>lon:([\d\.-]+)) (?P<heading>dir:([\d\.-]+)) (?P<speed>spd:([\d\.-]+)) (?P<lock>lck:([\d\.-]+)) (?P<time>time:([\d]+)) (?P<date>date:([\d]+)) (?P<status>trig:([\d]+))`)
	dataNames = dataRe.SubexpNames()
)

// UpdateShuttles send a request to iTrak API, gets updated shuttle info, and
// finally store updated records in db.
func (App *App) UpdateShuttles(dataFeed string, updateInterval int) {
	var st time.Duration
	for {
		// Sleep for n seconds before updating again
		log.Debugf("sleeping for %v", st)
		time.Sleep(st)
		if st == 0 {
			// Initialize the sleep timer after the first sleep.  This lets us sleep during errors
			// when we 'continue' back to the top of the loop without waiting to sleep for the first
			// update run.
			st = time.Duration(updateInterval) * time.Second
		}

		// Make request to our tracking data feed
		resp, err := http.Get(dataFeed)
		if err != nil {
			log.Errorf("error getting data feed: %v", err)
			continue
		}
		defer resp.Body.Close()

		// Read response body content
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("error reading data feed: %v", err)
			continue
		}

		delim := "eof"
		// Iterate through all vehicles returned by data feed
		vehiclesData := strings.Split(string(body), delim)
		// TODO: Figure out if this handles == 1 vehicle correctly or always assumes > 1.
		if len(vehiclesData) <= 1 {
			log.Warnf("found no vehicles delineated by '%s'", delim)
		}

		updated := 0
		for i := 0; i < len(vehiclesData)-1; i++ {
			match := dataRe.FindAllStringSubmatch(vehiclesData[i], -1)[0]

			// Store named capturing group and matching expression as a key value pair
			result := map[string]string{}
			for i, item := range match {
				result[dataNames[i]] = item
			}

			// Create new vehicle update & insert update into database
			update := VehicleUpdate{
				VehicleID: strings.Replace(result["id"], "Vehicle ID:", "", -1),
				Lat:       strings.Replace(result["lat"], "lat:", "", -1),
				Lng:       strings.Replace(result["lng"], "lon:", "", -1),
				Heading:   strings.Replace(result["heading"], "dir:", "", -1),
				Speed:     strings.Replace(result["speed"], "spd:", "", -1),
				Lock:      strings.Replace(result["lock"], "lck:", "", -1),
				Time:      strings.Replace(result["time"], "time:", "", -1),
				Date:      strings.Replace(result["date"], "date:", "", -1),
				Status:    strings.Replace(result["status"], "trig:", "", -1),
				Created:   time.Now()}

			if err := App.Updates.Insert(&update); err != nil {
				log.Errorf("error inserting vehicle update(%v): %v", update, err)
			} else {
				updated++
			}
		}
		log.Infof("sucessfully updated %d/%d vehicles", updated, len(vehiclesData)-1)
	}
}

// VehiclesHandler finds all the vehicles in the database.
func (App *App) VehiclesHandler(w http.ResponseWriter, r *http.Request) {
	// Find all vehicles in database
	var vehicles []Vehicle
	err := App.Vehicles.Find(bson.M{}).All(&vehicles)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each vehicle to client as JSON
	WriteJSON(w, vehicles)
}

// VehiclesCreateHandler adds a new vehicle to the database.
func (App *App) VehiclesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Create new vehicle object using request fields
	vehicle := Vehicle{}
	vehicleData := json.NewDecoder(r.Body)
	err := vehicleData.Decode(&vehicle)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Store new vehicle under vehicles collection
	err = App.Vehicles.Insert(&vehicle)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdatesHandler get the most recent update for each vehicle in the vehicles collection.
func (App *App) UpdatesHandler(w http.ResponseWriter, r *http.Request) {
	// Store updates for each vehicle
	var vehicles []Vehicle
	var updates []VehicleUpdate
	var update VehicleUpdate
	// Query all Vehicles
	err := App.Vehicles.Find(bson.M{}).All(&vehicles)
	// Handle errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Find recent updates for each vehicle
	for _, vehicle := range vehicles {
		err := App.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID}).Sort("-created").Limit(1).One(&update)

		if err == nil {
			updates = append(updates, update)
		}
	}
	// Convert updates to JSON
	WriteJSON(w, updates)
}
