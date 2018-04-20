package database

import (
	"time"
	"github.com/spf13/viper"
	"github.com/wtg/shuttletracker/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	//"github.com/wtg/shuttletracker/log"

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
	notifications 	*mgo.Collection

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
	db.notifications = db.session.DB("").C("notifications")


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
	route.ID = bson.NewObjectId().Hex()
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
	SetRouteActiveStatus(&route, time.Now())
	return route, err
}

// GetRoutes returns all Routes.
func (m *MongoDB) GetRoutes() ([]model.Route, error) {
	var routes []model.Route
	err := m.routes.Find(bson.M{}).All(&routes)
	for i := range routes {
		SetRouteActiveStatus(&routes[i], time.Now())
	}
	return routes, err
}

// ModifyRoute updates an existing Route by its ID.
func (m *MongoDB) ModifyRoute(route *model.Route) error {
	return m.routes.Update(bson.M{"id": route.ID}, route)
}

// CreateStop creates a Stop.
func (m *MongoDB) CreateStop(stop *model.Stop) error {
	stop.ID = bson.NewObjectId().Hex()
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

// GetStopsForRoute returns all stops on a specified route.
func (m *MongoDB) GetStopsForRoute(routeID string) ([]model.Stop, error){
	var stops []model.Stop

	err := m.stops.Find(bson.M{"routeId": routeID}).All(&stops)

	// Make sure stops on both routes are included.
	// Union has RouteID for East Route, but it's on both routes.
	// Not a very good solution to the problem -- FIXME

	// TODO: change RouteID for stops on multiple routes to be unique, empty,
	// or have all ID's for the routes they are on

	var stops_in_both bool = false
	for _, stop := range stops{
		if stop.Name == "Student Union"{
			stops_in_both = true
			break;
		}
	}

	if !stops_in_both{
		var union model.Stop
		err = m.stops.Find(bson.M{"name": "Student Union"}).One(&union)
		stops = append(stops, union)
	}

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

// AddMessage sets the current admin message
func (m *MongoDB) AddMessage(message *model.AdminMessage) error {
	message.ID = 1
	message.Created = time.Now()
	return m.messages.Insert(message)
}

// ClearMessage Clears the current message.
func (m *MongoDB) ClearMessage() error {
	message := model.AdminMessage{}
	message.ID = 1
	return m.messages.Remove(bson.M{"id": 1})
}

// GetCurrentMessage gets the most recent admin message
func (m *MongoDB) GetCurrentMessage() (model.AdminMessage, error) {
	message := model.AdminMessage{}
	err := m.messages.Find(bson.M{}).Sort("-created").One(&message)
	return message, err
}

// GetMessages gets the most recent admin messages
func (m *MongoDB) GetMessages() ([]model.AdminMessage, error) {
	messages := []model.AdminMessage{}
	err := m.messages.Find(bson.M{}).All(&messages)
	return messages, err
}

// UserExists tests if a given user exists in the admin database
func (m *MongoDB) UserExists(uname string) (bool, error) {
	query := m.users.Find(bson.M{"username": uname})
	n, err := query.Count()
	if n == 1 {
		return true, err
	}

	return false, err
}
// Creates a notification.
func (m *MongoDB) CreateNotification(notification *model.Notification) error {
	return m.notifications.Insert(&notification)
}

// Returns all notifications for the stop and route requested
func (m *MongoDB) GetNotificationsForStop(stopID string, routeID string) ([]model.Notification, error){
	var notifications []model.Notification
	err := m.notifications.Find(bson.M{"stop": stopID, "route": routeID}).All(&notifications)
	return notifications, err
}

// Deletes notifications based on stop and route
func (m *MongoDB) DeleteNotificationsForStop(stopID string, routeID string) (int, error){
	info, err := m.notifications.RemoveAll((bson.M{"stop": stopID, "route": routeID}))
	if err != nil {
		return 0, err
	}
	return info.Removed, nil
}
