package exporter

import (
	"encoding/json"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/model"
	"io/ioutil"

	"fmt"
	"github.com/wtg/shuttletracker/config"
	"github.com/wtg/shuttletracker/log"
	"os"
)

// Exporter contains the database (implements the database interface) to import/export from
type Exporter struct {
	db database.Database
}

// Dump represents all of the types in the database excluding updates
type Dump struct {
	Routes   []model.Route        `json:"routes"`
	Stops    []model.Stop         `json:"stops"`
	Vehicles []model.Vehicle      `json:"vehicles"`
	Users    []model.User         `json:"users"`
	Messages []model.AdminMessage `json:"messages"`
}

func (exp *Exporter) write(d Dump) {
	for _, r := range d.Routes {
		err := exp.db.CreateRoute(&r)
		if err != nil {
			log.WithError(err)
			return
		}
	}
	for _, r := range d.Stops {
		err := exp.db.CreateStop(&r)
		if err != nil {
			log.WithError(err)
			return
		}
	}
	for _, r := range d.Vehicles {
		err := exp.db.CreateVehicle(&r)
		if err != nil {
			log.WithError(err)
			return
		}
	}
	for _, r := range d.Users {
		err := exp.db.CreateUser(&r)
		if err != nil {
			log.WithError(err)
			return
		}
	}
	for _, r := range d.Messages {
		err := exp.db.AddMessage(&r)
		if err != nil {
			log.WithError(err)
			return
		}
	}
}

// Import imports the database from a given file, it assumes the database is empty
func (exp *Exporter) Import(file string) {
	d := Dump{}
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(raw, &d)
	exp.write(d)
}

// Export exports the database to the given file as a json
func (exp *Exporter) Export(dest string) {
	dump := Dump{}

	routes, err := exp.db.GetRoutes()
	if err != nil {
		log.WithError(err)
		return
	}
	dump.Routes = routes

	stops, err := exp.db.GetStops()
	if err != nil {
		log.WithError(err)
		return
	}
	dump.Stops = stops

	vehicles, err := exp.db.GetVehicles()
	if err != nil {
		log.WithError(err)
		return
	}
	dump.Vehicles = vehicles

	users, err := exp.db.GetUsers()
	if err != nil {
		log.WithError(err)
		return
	}
	dump.Users = users

	messages, err := exp.db.GetMessages()
	if err != nil {
		log.WithError(err)
		return
	}
	dump.Messages = messages

	out, err := json.Marshal(dump)
	if err != nil {
		log.WithError(err)
	}

	err = ioutil.WriteFile(dest, out, 0644)
	if err != nil {
		log.WithError(err)
	}
}
