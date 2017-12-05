package database

import (
	"errors"
	"time"

	"github.com/wtg/shuttletracker/model"
)

// Database is an interface that can be implemented by a database backend.
type Database interface {
	// Routes
	CreateRoute(route *model.Route) error
	DeleteRoute(routeID string) error
	GetRoute(routeID string) (model.Route, error)
	GetRoutes() ([]model.Route, error)
	ModifyRoute(route *model.Route) error

	// Stops
	CreateStop(stop *model.Stop) error
	DeleteStop(stopID string) error
	GetStops() ([]model.Stop, error)
	// GetStopsForRoute(routeID string) ([]model.Stop, error)
	// ModifyStop(stop *model.Stop) error

	// Vehicles
	CreateVehicle(vehicle *model.Vehicle) error
	DeleteVehicle(vehicleID string) error
	GetVehicle(vehicleID string) (model.Vehicle, error)
	GetVehicles() ([]model.Vehicle, error)
	GetEnabledVehicles() ([]model.Vehicle, error)
	ModifyVehicle(vehicle *model.Vehicle) error

	// Updates
	CreateUpdate(update *model.VehicleUpdate) error
	DeleteUpdatesBefore(before time.Time) (int, error)
	// GetUpdatesSince(since time.Time) ([]model.VehicleUpdate, error)
	GetUpdatesForVehicleSince(vehicleID string, since time.Time) ([]model.VehicleUpdate, error)
	GetLastUpdateForVehicle(vehicleID string) (model.VehicleUpdate, error)

	// Users
	GetUsers() ([]model.User, error)
}

var (
	ErrVehicleNotFound = errors.New("Vehicle not found.")
	ErrUpdateNotFound  = errors.New("Update not found.")
)
