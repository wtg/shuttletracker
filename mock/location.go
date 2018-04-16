package mock

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/wtg/shuttletracker"
)

// LocationService implements a mock of shuttletracker.LocationService.
type LocationService struct {
	mock.Mock
}

// CreateLocation creates a Location.
func (ls *LocationService) CreateLocation(location *shuttletracker.Location) error {
	args := ls.Called(location)
	return args.Error(0)
}

// DeleteLocationsBefore deletes Locations from before a certain time.
func (ls *LocationService) DeleteLocationsBefore(before time.Time) (int, error) {
	args := ls.Called(before)
	return args.Int(0), args.Error(1)
}

// LocationsSince gets Locations since a time for a certain Vehicle.
func (ls *LocationService) LocationsSince(vehicleID int, since time.Time) ([]*shuttletracker.Location, error) {
	args := ls.Called(vehicleID)
	return args.Get(0).([]*shuttletracker.Location), args.Error(1)
}

// LatestLocation returns the most recent Location for a Vehicle.
func (ls *LocationService) LatestLocation(vehicleID int) (*shuttletracker.Location, error) {
	args := ls.Called(vehicleID)
	return args.Get(0).(*shuttletracker.Location), args.Error(1)
}
