package tracking

import (
  "encoding/json"
  "net/http"
  "os"

  "gopkg.in/mgo.v2"
  log "github.com/Sirupsen/logrus"
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
}

func InitConfig() (*Configuration) {
  // Read app configuration file
  config, err := readConfiguration("conf.json")
  if err != nil {
    log.Fatalf("error reading configuration file: %v", err)
  }

  return config
}

func InitApp(Config *Configuration) (*App) {
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
