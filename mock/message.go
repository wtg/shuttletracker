package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/wtg/shuttletracker"
)

// MessageService implements a mock of shuttletracker.MessageService.
type MessageService struct {
	mock.Mock
}

// Message returns the Message.
func (ms *MessageService) Message() (*shuttletracker.Message, error) {
	args := ms.Called()
	return args.Get(0).(*shuttletracker.Message), args.Error(1)
}

// SetMessage updates the Message.
func (ms *MessageService) SetMessage(message *shuttletracker.Message) error {
	args := ms.Called(message)
	return args.Error(0)
}
