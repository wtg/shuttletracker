package mock

import (
	"github.com/stretchr/testify/mock"
)

// UserService implements a mock of shuttletracker.UserService.
type UserService struct {
	mock.Mock
}

// UserExists returns whether the User exists.
func (us *UserService) UserExists(username string) (bool, error) {
	args := us.Called(username)
	return args.Bool(0), args.Error(1)
}
