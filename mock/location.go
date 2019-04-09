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
func (ls *LocationService) LocationsSince(vehicleID int64, since time.Time) ([]*shuttletracker.Location, error) {
	args := ls.Called(vehicleID)
	return args.Get(0).([]*shuttletracker.Location), args.Error(1)
}

// LatestLocation returns the most recent Location for a Vehicle.
func (ls *LocationService) LatestLocation(vehicleID int64) (*shuttletracker.Location, error) {
	args := ls.Called(vehicleID)
	return args.Get(0).(*shuttletracker.Location), args.Error(1)
}

// LatestLocations returns the most recent Location for each Vehicle.
func (ls *LocationService) LatestLocations() ([]*shuttletracker.Location, error) {
	args := ls.Called()
	return args.Get(0).([]*shuttletracker.Location), args.Error(1)
}

// Location returns a Location by its ID.
func (ls *LocationService) Location(id int64) (*shuttletracker.Location, error) {
	args := ls.Called(id)
	return args.Get(0).(*shuttletracker.Location), args.Error(1)
}

// SubscribLocations returns a chan that receives each new Location.
func (ls *LocationService) SubscribeLocations() (chan *shuttletracker.Location) {
	args := ls.Called()
	return args.Get(0).(chan *shuttletracker.Location)
}
