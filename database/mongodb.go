package database

import (
	"time"

	"github.com/spf13/viper"
	"github.com/wtg/shuttletracker/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB implements Database with—you guessed it—MongoDB.
type MongoDB struct {
	session  *mgo.Session
	updates  *mgo.Collection
	vehicles *mgo.Collection
	routes   *mgo.Collection
	stops    *mgo.Collection
	users    *mgo.Collection
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

// CreateRoute creates a Route.
func (m *MongoDB) CreateRoute(route *model.Route) error {
	return m.routes.Insert(&route)
}

// DeleteRoute deletes a Route by its ID.
func (m *MongoDB) DeleteRoute(routeID string) error {
	return m.routes.Remove(bson.M{"id": routeID})
}

// GetRoute returns a Route by its ID.
func (m *MongoDB) GetRoute(routeID string) (model.Route, error) {
	var route model.Route
	err := m.routes.Find(bson.M{"id": routeID}).One(&route)
	return route, err
}

// GetRoutes returns all Routes.
func (m *MongoDB) GetRoutes() ([]model.Route, error) {
	var routes []model.Route
	err := m.routes.Find(bson.M{}).All(&routes)
	return routes, err
}

// ModifyRoute updates an existing Route by its ID.
func (m *MongoDB) ModifyRoute(route *model.Route) error {
	return m.routes.Update(bson.M{"id": route.ID}, route)
}

// CreateStop creates a Stop.
func (m *MongoDB) CreateStop(stop *model.Stop) error {
	return m.stops.Insert(&stop)
}

// DeleteStop deletes a Stop by its ID.
func (m *MongoDB) DeleteStop(stopID string) error {
	return m.stops.Remove(bson.M{"id": stopID})
}

// GetStop returns a Stop by its ID.
func (m *MongoDB) GetStop(stopID string) (model.Stop, error) {
	var stop model.Stop
	err := m.stops.Find(bson.M{"id": stopID}).One(&stop)
	return stop, err
}

// GetStops returns all Stops.
func (m *MongoDB) GetStops() ([]model.Stop, error) {
	var stops []model.Stop
	err := m.stops.Find(bson.M{}).All(&stops)
	return stops, err
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
func (m *MongoDB) GetLastUpdateForVehicle(vehicleID string) (model.VehicleUpdate, error) {
	var update model.VehicleUpdate
	err := m.updates.Find(bson.M{"vehicleID": vehicleID}).Sort("-created").One(&update)
	if err == mgo.ErrNotFound {
		return update, ErrUpdateNotFound
	}
	return update, err
}

// GetUpdatesForVehicleSince returns all updates since a time for a vehicle by its ID.
func (m *MongoDB) GetUpdatesForVehicleSince(vehicleID string, since time.Time) ([]model.VehicleUpdate, error) {
	var updates []model.VehicleUpdate
	err := m.updates.Find(bson.M{"vehicleID": vehicleID, "created": bson.M{"$gt": since}}).Sort("-created").All(&updates)
	return updates, err
}

// GetUsers returns all Users.
func (m *MongoDB) GetUsers() ([]model.User, error) {
	var users []model.User
	err := m.users.Find(bson.M{}).All(&users)
	return users, err
}

// CreateVehicle creates a Vehicle.
func (m *MongoDB) CreateVehicle(vehicle *model.Vehicle) error {
	return m.vehicles.Insert(&vehicle)
}

// DeleteVehicle deletes a Vehicle by its ID.
func (m *MongoDB) DeleteVehicle(vehicleID string) error {
	return m.vehicles.Remove(bson.M{"vehicleID": vehicleID})
}

// GetVehicle returns a Vehicle by its ID.
func (m *MongoDB) GetVehicle(vehicleID string) (model.Vehicle, error) {
	var vehicle model.Vehicle
	err := m.vehicles.Find(bson.M{"vehicleID": vehicleID}).One(&vehicle)
	if err == mgo.ErrNotFound {
		return vehicle, ErrVehicleNotFound
	}
	return vehicle, err
}

// GetVehicles returns all Vehicles.
func (m *MongoDB) GetVehicles() ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	err := m.vehicles.Find(bson.M{}).All(&vehicles)
	return vehicles, err
}

// GetEnabledVehicles returns all Vehicles that are enabled.
func (m *MongoDB) GetEnabledVehicles() ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	err := m.vehicles.Find(bson.M{"enabled": true}).All(&vehicles)
	return vehicles, err
}

// ModifyVehicle updates a Vehicle by its ID.
func (m *MongoDB) ModifyVehicle(vehicle *model.Vehicle) error {
	return m.vehicles.Update(bson.M{"vehicleID": vehicle.VehicleID}, vehicle)
}
