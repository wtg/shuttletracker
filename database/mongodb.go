package database

import (
	"time"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/wtg/shuttletracker/model"
)

// MongoDB implements Database with—you guessed it—MongoDB.
type MongoDB struct {
	session  *mgo.Session
	updates  *mgo.Collection
	vehicles *mgo.Collection
	routes   *mgo.Collection
	stops    *mgo.Collection
	users    *mgo.Collection
	messages *mgo.Collection
}

// MongoDBConfig contains information on how to connect to a MongoDB server.
type MongoDBConfig struct {
	MongoURL string
}

// NewMongoDB creates a MongoDB.
func NewMongoDB(cfg MongoDBConfig) (*MongoDB, error) {
	db := &MongoDB{}

	session, err := mgo.Dial(cfg.MongoURL)
	if err != nil {
		return nil, err
	}
	db.session = session

	db.updates = db.session.DB("").C("updates")
	db.vehicles = db.session.DB("").C("vehicles")
	db.routes = db.session.DB("").C("routes")
	db.stops = db.session.DB("").C("stops")
	db.users = db.session.DB("").C("users")
	db.messages = db.session.DB("").C("messages")

	// Ensure unique vehicle identification
	vehicleIndex := mgo.Index{
		Key:      []string{"vehicleID"},
		Unique:   true,
		DropDups: true}
	if err = db.vehicles.EnsureIndex(vehicleIndex); err != nil {
		return nil, err
	}

	// Create index on update vehicle ID and creation time to quickly find the most recent updates for specific vehicles.
	if err = db.updates.EnsureIndexKey("created"); err != nil {
		return nil, err
	}
	if err = db.updates.EnsureIndexKey("vehicleID"); err != nil {
		return nil, err
	}
	if err = db.updates.EnsureIndexKey("vehicleID", "created"); err != nil {
		return nil, err
	}

	// Index on enabled vehicles
	err = db.vehicles.EnsureIndexKey("enabled")

	return db, err
}

// NewMongoDBConfig creates a MongoDBConfig from a Viper instance.
func NewMongoDBConfig(v *viper.Viper) *MongoDBConfig {
	cfg := &MongoDBConfig{
		MongoURL: "localhost:27017",
	}
	v.SetDefault("database.mongourl", cfg.MongoURL)
	return cfg
}

// CreateUpdate creates an Update.
func (m *MongoDB) CreateUpdate(update *model.VehicleUpdate) error {
	return m.updates.Insert(&update)
}

// DeleteUpdatesBefore deletes all Updates that were created before a time.
func (m *MongoDB) DeleteUpdatesBefore(before time.Time) (int, error) {
	info, err := m.updates.RemoveAll(bson.M{"created": bson.M{"$lt": before}})
	if err != nil {
		return 0, err
	}
	return info.Removed, nil
}

// GetLastUpdateForVehicle returns the latest Update for a vehicle by its ID.
func (m *MongoDB) GetLastUpdateForVehicle(vehicleID int) (model.VehicleUpdate, error) {
	var update model.VehicleUpdate
	err := m.updates.Find(bson.M{"vehicleID": vehicleID}).Sort("-created").One(&update)
	if err == mgo.ErrNotFound {
		return update, ErrUpdateNotFound
	}
	return update, err
}

// GetUpdatesForVehicleSince returns all updates since a time for a vehicle by its ID.
func (m *MongoDB) GetUpdatesForVehicleSince(vehicleID int, since time.Time) ([]model.VehicleUpdate, error) {
	var updates []model.VehicleUpdate
	err := m.updates.Find(bson.M{"vehicleID": vehicleID, "created": bson.M{"$gt": since}}).Sort("-created").All(&updates)
	return updates, err
}
