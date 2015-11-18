package tracking

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

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
	Stops    *mgo.Collection
}

func InitConfig() *Configuration {
	// Read app configuration file
	config, err := readConfiguration("conf.json")
	if err != nil {
		log.Fatalf("error reading configuration file: %v", err)
	}

	return config
}

func InitApp(Config *Configuration) *App {
	// Connect to MongoDB
	session, err := mgo.Dial(Config.MongoURL + ":" + Config.MongoPort)
	if err != nil {
		log.Fatalf("mongoDB connection failed: %v", err)
	}
	// Create Shuttles object to store database session and collections
	app := App{
		session,
		session.DB("shuttle_tracking").C("updates"),
		session.DB("shuttle_tracking").C("vehicles"),
		session.DB("shuttle_tracking").C("routes"),
		session.DB("shuttle_tracking").C("stops"),
	}

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
	//error handling
	if err != nil {
		return err
	}
	//create a decoder for a file
	fileread := json.NewDecoder(file)
	//decode file into string

	//calling vehicles from vehicles.go
	Vehicles := []Vehicle{}

	//create map to hold variables
	var vehicles_map map[string][]map[string]interface{}

	//call decode on fileread to place items into map
	if err := fileread.Decode(&vehicles_map); err != nil {
		return err
	}

	for i := range vehicles_map["Vehicles"] {
		item := vehicles_map["Vehicles"][i]
		VehicleID, _ := item["VehicleID"].(string)
		VehicleName, _ := item["VehicleName"].(string)
		vehicle := Vehicle{VehicleID, VehicleName, time.Now(), time.Now()}
		Vehicles = append(Vehicles, vehicle)
	}

	for j := range Vehicles {
		err = app.Vehicles.Insert(&Vehicles[j])
		if err != nil {
			return err
		}
		fmt.Println("%v", Vehicles[j])
	}

	return nil
}

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
