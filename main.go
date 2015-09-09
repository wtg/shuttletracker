package main

import (
  "os"
  "fmt"
  "time"
  "strings"
  "net/http"
  "io/ioutil"
  "regexp"
  "encoding/json"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/mux"

  log "github.com/Sirupsen/logrus"
)

/**
 *  Shuttle Tracker 
 *   Auto Updater - send request to iTrak API,
 *                  get updated shuttle info,
 *                  store updated records in db
 */
type VehicleUpdate struct {
  Id          bson.ObjectId                 `bson:"_id,omitempty"`
  VehicleId   string     `json:"vehicleId"   bson:"vehicleId,omitempty"`
  Lat         string     `json:"lat"         bson:"lat"`
  Lng         string     `json:"lng"         bson:"lng"`
  Heading     string     `json:"heading"     bson:"heading"`
  Speed       string     `json:"speed"       bson:"speed"`
  Lock        string     `json:"lock"        bson:"lock"`
  Time        string     `json:"time"        bson:"time"`
  Date        string     `json:"date"        bson:"date"`
  Status      string     `json:"status"      bson:"status"`
  Created     time.Time  `json:"created"     bson:"created"`
}

var (
  // Match each API field with any number (+)
  //   of the previous expressions (\d digit, \. escaped period, - negative number)
  //   Specify named capturing groups to store each field from data feed
  dataRe = regexp.MustCompile(`(?P<id>Vehicle ID:([\d\.]+)) (?P<lat>lat:([\d\.-]+)) (?P<lng>lon:([\d\.-]+)) (?P<heading>dir:([\d\.-]+)) (?P<speed>spd:([\d\.-]+)) (?P<lock>lck:([\d\.-]+)) (?P<time>time:([\d]+)) (?P<date>date:([\d]+)) (?P<status>trig:([\d]+))`)
  dataNames = dataRe.SubexpNames()
)

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
      continue;
    }
    defer resp.Body.Close()

    // Read response body content
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      log.Errorf("error reading data feed: %v", err)
      continue;
    }

    delim := "eof"
    // Iterate through all vehicles returned by data feed
    vehicles_data := strings.Split(string(body), delim)
    // TODO: Figure out if this handles == 1 vehicle correctly or always assumes > 1.
    if len(vehicles_data) <= 1 {
      log.Warnf("found no vehicles delineated by '%s'", delim)
    }

    updated := 0
    for i := 1; i < len(vehicles_data)-1; i++ {
      match := dataRe.FindAllStringSubmatch(vehicles_data[i], -1)[0]

      // Store named capturing group and matching expression as a key value pair
      result := map[string]string{}
      for i, item := range match {
        result[dataNames[i]] = item
      }

      // Create new vehicle update & insert update into database
      update := VehicleUpdate { 
        Id:        bson.NewObjectId(),
        VehicleId: strings.Replace(result["id"], "Vehicle ID:", "", -1),
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
    log.Infof("sucessfully updated %d/%d vehicles", updated, len(vehicles_data)-1)
  }
}

/**
 *  Route handlers - API requests,
 *                   serve view files
 */

type Vehicle struct {
  Id          bson.ObjectId                 `bson:"_id,omitempty"`
  VehicleId   string     `json:"vehicleId"   bson:"vehicleId,omitempty"`
  VehicleName string     `json:"vehicleName" bson:"vehicleName"`
  Created     time.Time  `json:"created"     bson:"created"`
}

/**
 * Find all vehicles in the database 
 *
 */
func (App *App) VehiclesHandler(w http.ResponseWriter, r *http.Request) {
  // Find all vehicles in database
  var vehicles []Vehicle
  err := App.Vehicles.Find(bson.M{}).All(&vehicles)
  // Handle query errors
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  // Send each vehicle to client as JSON
  vehiclesJSON, err := json.MarshalIndent(vehicles, "", " ")
  fmt.Fprintf(w, string(vehiclesJSON))
}

/**
 * Add a new vehicle to the database 
 *
 */
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

/**
 * Vehicle Updates - Get most recent update for each
 *                   vehicle in the vehicles collection
 */
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
  for _,vehicle := range vehicles {
    err := App.Updates.Find(bson.M{"vehicleId": vehicle.VehicleId}).Sort("-created").Limit(1).One(&update)

    if err == nil {
      updates = append(updates, update)
    }
  }
  // Convert updates to JSON
  u, err := json.MarshalIndent(updates, "", " ")
  // Handle JSON parsing errors
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  // Send updates to client 
  fmt.Fprint(w, string(u))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "index.html")
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "admin.html")
}

/**
 *  Main - connect to database, 
 *         handle routing,
 *         start tracker go routine,
 *         serve requests
 */

type Configuration struct {
  DataFeed        string
  UpdateInterval  int
  MongoUrl        string
  MongoPort       string
}

type App struct {
  Session    *mgo.Session
  Updates    *mgo.Collection
  Vehicles   *mgo.Collection
}

func ReadConfiguration(fileName string) (*Configuration, error) {
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

func main() {
  // Read app configuration file 
  config, err := ReadConfiguration("conf.json")
  if err != nil {
    log.Fatalf("error reading configuration file: %v", err)
  }

  // Connect to MongoDB
  session, err := mgo.Dial(config.MongoUrl + ":" + config.MongoPort)
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
  }

  // Start auto updater 
  go App.UpdateShuttles(config.DataFeed, config.UpdateInterval)

  // Routing 
  r := mux.NewRouter()
  r.HandleFunc("/", IndexHandler).Methods("GET")
  r.HandleFunc("/admin", AdminHandler).Methods("GET")
  r.HandleFunc("/admin/vehicles", AdminHandler).Methods("GET")
  r.HandleFunc("/admin/tracking", AdminHandler).Methods("GET")
  r.HandleFunc("/admin/stops", AdminHandler).Methods("GET")
  r.HandleFunc("/vehicles", App.VehiclesHandler).Methods("GET")
  r.HandleFunc("/vehicles/create", App.VehiclesCreateHandler).Methods("POST")
  r.HandleFunc("/updates", App.UpdatesHandler).Methods("GET")
  // Static files
  r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components/"))))
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
  // Serve requests
  http.Handle("/", r)
  if err := http.ListenAndServe(":8080", r); err != nil {
    log.Fatalf("Unable to ListenAndServe: %v", err)
  }
}
