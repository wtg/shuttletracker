package database

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/wtg/shuttletracker/model"
)

// Mock implements the database interface. Great for tests.
type Mock struct {
	mock.Mock
}

// CreateRoute creates a route.
func (db *Mock) CreateRoute(route *model.Route) error {
	args := db.Called(route)
	return args.Error(0)
}

// DeleteRoute deletes a route.
func (db *Mock) DeleteRoute(routeID string) error {
	args := db.Called(routeID)
	return args.Error(0)
}

// GetRoute gets a route.
func (db *Mock) GetRoute(routeID string) (model.Route, error) {
	args := db.Called(routeID)
	return args.Get(0).(model.Route), args.Error(1)
}

// GetRoutes gets all routes.
func (db *Mock) GetRoutes() ([]model.Route, error) {
	args := db.Called()
	return args.Get(0).([]model.Route), args.Error(1)
}

// ModifyRoute modifies a route.
func (db *Mock) ModifyRoute(route *model.Route) error {
	args := db.Called(route)
	return args.Error(0)
}

// CreateStop creates a stop.
func (db *Mock) CreateStop(stop *model.Stop) error {
	args := db.Called(stop)
	return args.Error(0)
}

// DeleteStop deletes a stop.
func (db *Mock) DeleteStop(stopID string) error {
	args := db.Called(stopID)
	return args.Error(0)
}

// GetStops gets all stops.
func (db *Mock) GetStops() ([]model.Stop, error) {
	args := db.Called()
	return args.Get(0).([]model.Stop), args.Error(1)
}

// CreateVehicle creates a vehicle.
func (db *Mock) CreateVehicle(vehicle *model.Vehicle) error {
	args := db.Called(vehicle)
	return args.Error(0)
}

// DeleteVehicle deletes a vehicle.
func (db *Mock) DeleteVehicle(vehicleID string) error {
	args := db.Called(vehicleID)
	return args.Error(0)
}

// GetVehicle gets a vehicle.
func (db *Mock) GetVehicle(vehicleID string) (model.Vehicle, error) {
	args := db.Called(vehicleID)
	return args.Get(0).(model.Vehicle), args.Error(1)
}

// GetVehicles gets all vehicles.
func (db *Mock) GetVehicles() ([]model.Vehicle, error) {
	args := db.Called()
	return args.Get(0).([]model.Vehicle), args.Error(1)
}

// GetEnabledVehicles gets all enabled vehicles.
func (db *Mock) GetEnabledVehicles() ([]model.Vehicle, error) {
	args := db.Called()
	return args.Get(0).([]model.Vehicle), args.Error(1)
}

// ModifyVehicle modifies a vehicle.
func (db *Mock) ModifyVehicle(vehicle *model.Vehicle) error {
	args := db.Called(vehicle)
	return args.Error(0)
}

// CreateUpdate creates an update.
func (db *Mock) CreateUpdate(update *model.VehicleUpdate) error {
	args := db.Called(update)
	return args.Error(0)
}

// DeleteUpdatesBefore deletes all updates before a time.
func (db *Mock) DeleteUpdatesBefore(before time.Time) (int, error) {
	args := db.Called(before)
	return args.Int(0), args.Error(1)
}

// GetUpdatesForVehicleSince gets all updates for a vehicle since a time.
func (db *Mock) GetUpdatesForVehicleSince(vehicleID string, since time.Time) ([]model.VehicleUpdate, error) {
	args := db.Called(vehicleID, since)
	return args.Get(0).([]model.VehicleUpdate), args.Error(1)
}

// GetLastUpdateForVehicle gets the most recent update for a vehicle.
func (db *Mock) GetLastUpdateForVehicle(vehicleID string) (model.VehicleUpdate, error) {
	args := db.Called(vehicleID)
	return args.Get(0).(model.VehicleUpdate), args.Error(1)
}

// GetUsers gets all users.
func (db *Mock) GetUsers() ([]model.User, error) {
	args := db.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

// UserExists gets all users.
func (db *Mock) UserExists() (bool, error) {
	args := db.Called()
	return args.Get(0).(bool), args.Error(1)
}

// AddMessage adds a message.
func (db *Mock) AddMessage(message *model.AdminMessage) error {
	args := db.Called(message)
	return args.Error(0)
}

// GetCurrentMessage gets the current message.
func (db *Mock) GetCurrentMessage() (model.AdminMessage, error) {
	args := db.Called()
	return args.Get(0).(model.AdminMessage), args.Error(1)
}

// GetMessages gets all messages.
func (db *Mock) GetMessages() ([]model.AdminMessage, error) {
	args := db.Called()
	return args.Get(0).([]model.AdminMessage), args.Error(1)
}

// ClearMessage clears the message.
func (db *Mock) ClearMessage() error {
	args := db.Called()
	return args.Error(0)
}
