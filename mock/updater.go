package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/wtg/shuttletracker"
)

// UpdaterService implements a mock of shuttletracker.UpdaterService.
type UpdaterService struct {
	mock.Mock
}

// GetLastResponse returns the most recent iTRAK datafeed response.
func (us *UpdaterService) GetLastResponse() *shuttletracker.DataFeedResponse {
	args := us.Called()
	return args.Get(0).(*shuttletracker.DataFeedResponse)
}
