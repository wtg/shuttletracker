package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/wtg/shuttletracker"
)

// StopService implements a mock of shuttletracker.StopService.
type StopService struct {
	mock.Mock
}

// CreateStop creates a Stop.
func (ss *StopService) CreateStop(stop *shuttletracker.Stop) error {
	args := ss.Called(stop)
	return args.Error(0)
}

// DeleteStop deletes a Stop.
func (ss *StopService) DeleteStop(id int64) error {
	args := ss.Called(id)
	return args.Error(0)
}

// Stops gets all stops.
func (ss *StopService) Stops() ([]*shuttletracker.Stop, error) {
	args := ss.Called()
	return args.Get(0).([]*shuttletracker.Stop), args.Error(1)
}

// Stop gets a Stop.
func (ss *StopService) Stop(stopID int64) (*shuttletracker.Stop, error) {
	args := ss.Called(stopID)
	return args.Get(0).(*shuttletracker.Stop), args.Error(1)
}
