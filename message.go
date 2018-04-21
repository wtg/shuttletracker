package shuttletracker

import (
	"time"
)

// Message represents a message displayed to users.
type Message struct {
	ID      int       `json:"id"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
	Display bool      `json:"display"`
	Created time.Time `json:"created"`
}

// MessageService is an interface for interacting with Messages.
type MessageService interface {
	Message() (*Message, error)
	SetMessage(message *Message) error
}
