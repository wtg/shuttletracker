package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/wtg/shuttletracker"
)

// UpdaterService implements a mock of shuttletracker.UpdaterService.
type UpdaterService struct {
	mock.Mock
}

func (us *UpdaterService) GetLastResponse() *shuttletracker.DataFeedResponse {
	args := us.Called()
	return args.Get(0).(*shuttletracker.DataFeedResponse)
}
