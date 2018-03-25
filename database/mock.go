package database

import (
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/wtg/shuttletracker/model"
)

type Mock struct {
	mock.Mock
}

func (db *Mock) CreateRoute(route *model.Route) error {
	args := db.Called(route)
	return args.Error(0)
}

func (db *Mock) DeleteRoute(routeID string) error {
	args := db.Called(routeID)
	return args.Error(0)
}

func (db *Mock) GetRoute(routeID string) (model.Route, error) {
	args := db.Called(routeID)
	return args.Get(0).(model.Route), args.Error(1)
}

func (db *Mock) GetRoutes() ([]model.Route, error) {
	args := db.Called()
	return args.Get(0).([]model.Route), args.Error(1)
}

func (db *Mock) ModifyRoute(route *model.Route) error {
	args := db.Called(route)
	return args.Error(0)
}

func (db *Mock) CreateStop(stop *model.Stop) error {
	args := db.Called(stop)
	return args.Error(0)
}

func (db *Mock) DeleteStop(stopID string) error {
	args := db.Called(stopID)
	return args.Error(0)
}

func (db *Mock) GetStops() ([]model.Stop, error) {
	args := db.Called()
	return args.Get(0).([]model.Stop), args.Error(1)
}

func (db *Mock) CreateVehicle(vehicle *model.Vehicle) error {
	args := db.Called(vehicle)
	return args.Error(0)
}

func (db *Mock) DeleteVehicle(vehicleID string) error {
	args := db.Called(vehicleID)
	return args.Error(0)
}

func (db *Mock) GetVehicle(vehicleID string) (model.Vehicle, error) {
	args := db.Called(vehicleID)
	return args.Get(0).(model.Vehicle), args.Error(1)
}

func (db *Mock) GetVehicles() ([]model.Vehicle, error) {
	args := db.Called()
	return args.Get(0).([]model.Vehicle), args.Error(1)
}

func (db *Mock) GetEnabledVehicles() ([]model.Vehicle, error) {
	args := db.Called()
	return args.Get(0).([]model.Vehicle), args.Error(1)
}

func (db *Mock) ModifyVehicle(vehicle *model.Vehicle) error {
	args := db.Called(vehicle)
	return args.Error(0)
}

func (db *Mock) CreateUpdate(update *model.VehicleUpdate) error {
	args := db.Called(update)
	return args.Error(0)
}

func (db *Mock) DeleteUpdatesBefore(before time.Time) (int, error) {
	args := db.Called(before)
	return args.Int(0), args.Error(1)
}

func (db *Mock) GetUpdatesForVehicleSince(vehicleID string, since time.Time) ([]model.VehicleUpdate, error) {
	args := db.Called(vehicleID, since)
	return args.Get(0).([]model.VehicleUpdate), args.Error(1)
}

func (db *Mock) GetLastUpdateForVehicle(vehicleID string) (model.VehicleUpdate, error) {
	args := db.Called(vehicleID)
	return args.Get(0).(model.VehicleUpdate), args.Error(1)
}

func (db *Mock) GetUsers() ([]model.User, error) {
	args := db.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (db *Mock) AddMessage(message *model.AdminMessage) error {
	args := db.Called(message)
	return args.Error(0)
}

func (db *Mock) GetCurrentMessage() (model.AdminMessage, error) {
	args := db.Called()
	return args.Get(0).(model.AdminMessage), args.Error(1)
}

func (db *Mock) GetMessages() ([]model.AdminMessage, error) {
	args := db.Called()
	return args.Get(0).([]model.AdminMessage), args.Error(1)
}

func (db *Mock) ClearMessage() error {
	args := db.Called()
	return args.Error(0)
}
