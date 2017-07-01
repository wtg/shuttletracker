package database

import (
	"gopkg.in/mgo.v2"
	"github.com/wtg/shuttletracker/log"
)

type Database struct {
	session *mgo.Session
	Updates *mgo.Collection
	Vehicles *mgo.Collection
Routes   *mgo.Collection
Stops    *mgo.Collection
Users    *mgo.Collection
}

type Config struct {
	MongoURL string
}

func New(cfg Config) *Database {
	db := &Database{}

	session, err := mgo.Dial(cfg.MongoURL)
	if err != nil {
		log.Errorf("MongoDB connection to \"%v\" failed: %v", cfg.MongoURL, err)
	}
	db.session = session

	db.Updates = db.session.DB("").C("updates")
	db.Vehicles = db.session.DB("").C("vehicles")
	db.Routes = db.session.DB("").C("routes")
	db.Stops = db.session.DB("").C("stops")
	db.Users = db.session.DB("").C("users")

	// Ensure unique vehicle identification
	vehicleIndex := mgo.Index{
		Key:      []string{"vehicleID"},
		Unique:   true,
		DropDups: true}
	db.Vehicles.EnsureIndex(vehicleIndex)

	// Create index on update created time to quickly find the most recent updates
	db.Updates.EnsureIndexKey("created")

	return db
}
