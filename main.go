package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
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
		for i := 1; i < len(vehiclesData)-1; i++ {
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
	writeJSON(w, vehicles)
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
	writeJSON(w, updates)
}

// Coord objects contain the lat/long coordinates to draw routes
type Coord struct {
	Lat     float64       `bson:"lat"`
	Lng     float64       `bson:"lng"`
	RouteID bson.ObjectId `bson:routeID"`
	Created time.Time     `bson:created"`
	Updated time.Time     `bson:updated"`
}

// Route represents a set of coordinates to draw a path on our tracking map
type Route struct {
	Name        string    `json:"name"           bson:"name"`
	Description string    `json:"description"    bson:"description"`
	StartTime   string    `json:"startTime"      bson:"startTime"`
	EndTime     string    `json:"endTime"        bson:"endTime"`
	Enabled     bool      `json:"enabled,string" bson:"enabled"`
	Color       string    `json:"color"          bson:"color"`
	Width       int       `json:"width,string"   bson:"width"`
	Created     time.Time `json:"created"        bson:"created"`
	Updated     time.Time `json:"updated"        bson:"updated"`
}

// Stop indicates where a tracked object is scheduled to arrive
type Stop struct {
	Name        string    `json:"name"        bson:"name"`
	Phonetic    string    `json:"phonetic"    bson:"phonetic"`
	Description string    `json:"description" bson:"description"`
	Address     string    `json:"address"     bson:"address"`
	TimeServed  string    `json:"timeServed"  bson:"timeServed"`
	Lat         float64   `json:"lat"         bson:"lat"`
	Lng         float64   `json:"lng"         bson:"lng"`
	Enabled     bool      `json:"enabled"     bson:"enabled"`
	Created     time.Time `json:"created"     bson:"created"`
	Updated     time.Time `json:"updated"     bson:"updated"`
}

// RouteStop allows stops to be placed on one or more routes
type RouteStop struct {
	StopID  bson.ObjectId `bson:"stopID",omitempty"`
	RouteID bson.ObjectId `bson:"routeID",omitempty"`
}

// RoutesHandler finds all of the routes in the database
func (App *App) RoutesHandler(w http.ResponseWriter, r *http.Request) {
	// Find all routes in database
	var routes []Route
	err := App.Routes.Find(bson.M{}).All(&routes)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each route to client as JSON
	writeJSON(w, routes)
}

// RoutesCreateHandler adds a new route to the database
func (App *App) RoutesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new route object using request fields
	route := Route{}
	routeData := json.NewDecoder(r.Body)
	err := routeData.Decode(&route)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Store new route under routes collection
	err = App.Routes.Insert(&route)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// IndexHandler serves the index page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// AdminHandler serves the admin page.
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "admin.html")
}

// Configuration holds the settings for connecting to outside resources.
type Configuration struct {
	DataFeed       string
	UpdateInterval int
	MongoURL       string
	MongoPort      string
}

// App holds references to Mongo resources.
type App struct {
	Session  *mgo.Session
	Updates  *mgo.Collection
	Vehicles *mgo.Collection
	Routes   *mgo.Collection
}

func readConfiguration(fileName string) (*Configuration, error) {
	// Open config file and decode JSON to Configuration struct
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	config := Configuration{}
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func writeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return nil
}

func main() {
	// Read app configuration file
	config, err := readConfiguration("conf.json")
	if err != nil {
		log.Fatalf("error reading configuration file: %v", err)
	}

	// Connect to MongoDB
	session, err := mgo.Dial(config.MongoURL + ":" + config.MongoPort)
	if err != nil {
		log.Fatalf("mongoDB connection failed: %v", err)
	}
	// close Mongo session when server terminates
	defer session.Close()

	// Create Shuttles object to store database session and collections
	App := &App{
		session,
		session.DB("shuttle_tracking").C("updates"),
		session.DB("shuttle_tracking").C("vehicles"),
		session.DB("shuttle_tracking").C("routes"),
	}

	// Start auto updater
	go App.UpdateShuttles(config.DataFeed, config.UpdateInterval)

	// Routing
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/admin", AdminHandler).Methods("GET")
	r.HandleFunc("/admin/{*}", AdminHandler).Methods("GET")
	r.HandleFunc("/vehicles", App.VehiclesHandler).Methods("GET")
	r.HandleFunc("/vehicles/create", App.VehiclesCreateHandler).Methods("POST")
	r.HandleFunc("/updates", App.UpdatesHandler).Methods("GET")
	r.HandleFunc("/routes", App.RoutesHandler).Methods("GET")
	r.HandleFunc("/routes/create", App.RoutesCreateHandler).Methods("POST")
	// Static files
	r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components/"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	// Serve requests
	http.Handle("/", r)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Unable to ListenAndServe: %v", err)
	}
}
