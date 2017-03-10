package tracking

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// Configuration holds the settings for connecting to outside resources.
type Configuration struct {
	DataFeed             string
	UpdateInterval       int
	MongoURL             string
	MongoPort            string
	GoogleMapAPIKey      string
	GoogleMapMinDistance int
}

// App holds references to Mongo resources.
type App struct {
	Config   *Configuration
	Session  *mgo.Session
	Updates  *mgo.Collection
	Vehicles *mgo.Collection
	Routes   *mgo.Collection
	Stops    *mgo.Collection
}

// InitConfig loads and return the app config.
func InitConfig() *Configuration {
	// Read app configuration file
	config, err := readConfiguration("conf.json")
	if err != nil {
		log.Fatalf("error reading configuration file: %v", err)
	}

	return config
}

// InitApp initializes the application given a config and connects to backends.
// It also seeds any needed information to the database.
func InitApp(Config *Configuration) *App {
	// Connect to MongoDB
	session, err := mgo.Dial(Config.MongoURL + ":" + Config.MongoPort)
	if err != nil {
		log.Fatalf("mongoDB connection failed: %v", err)
	}
	// Create Shuttles object to store database session and collections
	app := App{
		Config,
		session,
		session.DB("shuttle_tracking").C("updates"),
		session.DB("shuttle_tracking").C("vehicles"),
		session.DB("shuttle_tracking").C("routes"),
		session.DB("shuttle_tracking").C("stops"),
	}

	// Ensure unique vehicle identification
	vehicleIndex := mgo.Index{
		Key:      []string{"vehicleID"},
		Unique:   true,
		DropDups: true}
	app.Vehicles.EnsureIndex(vehicleIndex)

	// Read vehicle configuration file
	serr := readSeedConfiguration("seed/vehicle_seed.json", &app)
	if serr != nil {
		log.Fatalf("error reading vehicle configuration file: %v", serr)
	}
	return &app
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

//readSeedConfiguration adds a new vehicle to the database from seed.
func readSeedConfiguration(fileName string, app *App) error {
	// Open seed_vehicle config file and decode JSON to app struct
	file, err := os.Open(fileName)

	// Error handling
	if err != nil {
		log.Warn(err)
	}
	// Create a decoder for a file
	fileread := json.NewDecoder(file)

	// Create map for json data and slice for vehicles
	var vehiclesMap map[string][]map[string]interface{} // map with string as key and ,list of map with string as key and anything as value, as value
	Vehicles := []Vehicle{}                             // list of default vehicle object

	// Call decode on fileread to place items into map
	if err := fileread.Decode(&vehiclesMap); err != nil {
		log.Warn(err)
	}

	// Initialize our vehicles
	for i := range vehiclesMap["Vehicles"] {
		item := vehiclesMap["Vehicles"][i]
		VehicleID, _ := item["VehicleID"].(string)
		VehicleName, _ := item["VehicleName"].(string)
		vehicle := Vehicle{VehicleID, VehicleName, time.Now(), time.Now()}
		Vehicles = append(Vehicles, vehicle)
	}

	// Add vehicles to the database
	for j := range Vehicles {
		app.Vehicles.Insert(&Vehicles[j])
	}

	return nil
}

// WriteJSON writes the data as JSON.
func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return nil
}
