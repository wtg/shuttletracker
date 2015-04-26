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

func (Shuttles *Shuttles) UpdateShuttles(dataFeed string, updateInterval int) {
  for {
    // Reference updates collection and close db session upon exit
    UpdatesCollection := Shuttles.Session.DB("shuttle_tracking").C("updates")

    // Make request to our tracking data feed
    resp, err := http.Get(dataFeed)
    if err != nil {
      continue;
    }
    defer resp.Body.Close()

    // Read response body content
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      continue;
    }

    delim := "eof"
    // Iterate through all vehicles returned by data feed
    vehicles_data := strings.Split(string(body), delim)
    for i := 1; i < len(vehicles_data)-1; i++ {

      // Match eatch API field with any number (+)
      //   of the previous expressions (\d digit, \. escaped period, - negative number)
      //   Specify named capturing groups to store each field from data feed
      re := regexp.MustCompile(`(?P<id>Vehicle ID:([\d\.]+)) (?P<lat>lat:([\d\.-]+)) (?P<lng>lon:([\d\.-]+)) (?P<heading>dir:([\d\.-]+)) (?P<speed>spd:([\d\.-]+)) (?P<lock>lck:([\d\.-]+)) (?P<time>time:([\d]+)) (?P<date>date:([\d]+)) (?P<status>trig:([\d]+))`)
      n := re.SubexpNames()
      match := re.FindAllStringSubmatch(vehicles_data[i], -1)[0]

      // Store named capturing group and matching expression as a key value pair
      result := map[string]string{}
      for i, item := range match {
        result[n[i]] = item
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

      err := UpdatesCollection.Insert(&update)

      if err != nil {
        fmt.Println(err.Error())
      }
    }

    // Sleep for n seconds before updating again
    time.Sleep(time.Duration(updateInterval) * time.Second)
  }
}

/**
 *  Route handlers - API requests,
 *                   serve view files
 */

func IndexHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "index.html")
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "dashboard.html")
}

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
func (Shuttles *Shuttles) VehiclesHandler(w http.ResponseWriter, r *http.Request) {
  // Find all vehicles in database
  var vehicles []Vehicle
  VehiclesCollection := Shuttles.Session.DB("shuttle_tracking").C("vehicles")
  err := VehiclesCollection.Find(bson.M{}).All(&vehicles)
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
func (Shuttles *Shuttles) VehiclesCreateHandler(w http.ResponseWriter, r *http.Request) {
  // Create new vehicle object using request fields
  vehicle := Vehicle{}
  vehicleData := json.NewDecoder(r.Body)
  err := vehicleData.Decode(&vehicle)
  // Error handling
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  // Store new vehicle under vehicles collection
  VehiclesCollection := Shuttles.Session.DB("shuttle_tracking").C("vehicles")
  err = VehiclesCollection.Insert(&vehicle)
  // Error handling
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

/**
 * Vehicle Updates - Get most recent update for each
 *                   vehicle in the vehicles collection
 */
func (Shuttles *Shuttles) UpdatesHandler(w http.ResponseWriter, r *http.Request) {
  // Access vehicles and updates collections in shuttle tracking database 
  UpdatesCollection := Shuttles.Session.DB("shuttle_tracking").C("updates")
  VehiclesCollection := Shuttles.Session.DB("shuttle_tracking").C("vehicles")
  // Store updates for each vehicle
  var vehicles []Vehicle
  var updates []VehicleUpdate
  var update VehicleUpdate
  // Query all Vehicles 
  err := VehiclesCollection.Find(bson.M{}).All(&vehicles)
  // Handle errors
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  // Find recent updates for each vehicle
  for _,vehicle := range vehicles {
    err := UpdatesCollection.Find(bson.M{"vehicleId": vehicle.VehicleId}).Sort("-created").Limit(1).One(&update)
    updates = append(updates, update)

    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
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

type Shuttles struct {
  Session *mgo.Session
}

func ReadConfiguration(fileName string) Configuration { 
  // Open config file and decode JSON to Configuration struct
  file, _ := os.Open(fileName)
  decoder := json.NewDecoder(file)
  config := Configuration{}
  err := decoder.Decode(&config)
  if err != nil {
    fmt.Println("Unable to read config file: ")
    os.Exit(1)
  }
  return config
}

func main() {
  // Read app configuration file 
  config := ReadConfiguration("conf.json")

  // Connect to MongoDB
  session, err := mgo.Dial(config.MongoUrl + ":" + config.MongoPort)
  if err != nil {
    fmt.Println("MongoDB connection failed")
    os.Exit(1)
  }
  // close Mongo session when server terminates
  defer session.Close()

  // Create Shuttles object to store database session information
  Shuttles := &Shuttles{session}

  // Start auto updater 
  go Shuttles.UpdateShuttles(config.DataFeed, config.UpdateInterval)

  // Routing 
  r := mux.NewRouter()
  r.HandleFunc("/", IndexHandler).Methods("GET")
  r.HandleFunc("/admin", AdminHandler).Methods("GET")
  r.HandleFunc("/admin/vehicles", AdminHandler).Methods("GET")
  r.HandleFunc("/admin/tracking", AdminHandler).Methods("GET")
  r.HandleFunc("/vehicles", Shuttles.VehiclesHandler).Methods("GET")
  r.HandleFunc("/vehicles/create", Shuttles.VehiclesCreateHandler).Methods("POST")
  r.HandleFunc("/updates", Shuttles.UpdatesHandler).Methods("GET")
  // Static files
  r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components/"))))
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
  // Serve requests
  http.Handle("/", r)
  http.ListenAndServe(":8080", r)
}