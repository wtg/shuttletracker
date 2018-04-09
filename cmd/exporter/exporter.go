package exporter

import (
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/model"
  "io/ioutil"
  "encoding/json"

	"github.com/wtg/shuttletracker/config"
	"github.com/wtg/shuttletracker/log"
)

type Exporter struct {
	db database.Database
}

type Dump struct {
	Routes   []model.Route        `json:"routes"`
	Stops    []model.Stop         `json:"stops"`
	Vehicles []model.Vehicle      `json:"vehicles"`
	Users    []model.User         `json:users`
	Messages []model.AdminMessage `json:messages`
}

func (exp *Exporter) Export(dest string) {
	dump := Dump{}

	routes, err := exp.db.GetRoutes()
  if(err != nil){
    log.WithError(err)
    return
  }
  dump.Routes = routes

  stops, err := exp.db.GetStops()
  if(err != nil){
    log.WithError(err)
    return
  }
  dump.Stops = stops

  vehicles, err := exp.db.GetVehicles()
  if(err != nil){
    log.WithError(err)
    return
  }
  dump.Vehicles = vehicles

  users, err := exp.db.GetUsers()
  if(err != nil){
    log.WithError(err)
    return
  }
  dump.Users = users

  messages, err := exp.db.GetMessages()
  if(err != nil){
    log.WithError(err)
    return
  }
  dump.Messages = messages

  out, err := json.Marshal(dump)
  if(err != nil){
    log.WithError(err)
  }

  err  = ioutil.WriteFile(dest, out, 0644)
  if(err != nil){
    log.WithError(err)
  }
}
// 
// func Main() {
// 	log.Info("Shuttle Tracker starting...")
//
// 	// Config
// 	cfg, err := config.New()
// 	if err != nil {
// 		log.WithError(err).Error("Could not create config.")
// 		return
// 	}
// 	// Database
// 	db, err := database.NewMongoDB(*cfg.Database)
// 	if err != nil {
// 		log.WithError(err).Errorf("MongoDB connection to \"%v\" failed.", cfg.Database.MongoURL)
// 		return
// 	}
// 	e := Exporter{db}
// 	e.Export("dump.json")
//
// }
