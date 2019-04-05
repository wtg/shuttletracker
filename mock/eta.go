package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/wtg/shuttletracker"
)

// ETAService implements a mock of shuttletracker.ETAService.
type ETAService struct {
	mock.Mock
}

// CurrentETAs returns the ETA service's ETAs.
func (es *ETAService) CurrentETAs() map[int64]shuttletracker.VehicleETA {
	args := es.Called()
	return args.Get(0).(map[int64]shuttletracker.VehicleETA)
}

// Subscribe allows callers to provide a callback to receive VehicleETAs.
func (es *ETAService) Subscribe(f func(shuttletracker.VehicleETA)) {
	es.Called(f)
}
