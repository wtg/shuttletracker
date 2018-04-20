package shuttletracker

import (
	"time"
)

// AdminMessage represents a message popup for the user from the site administrator
type Message struct {
	ID      int       `json:"id"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
	Display bool      `json:"display"`
	Created time.Time `json:"created"`
}

type MessageService interface {
	Message() (*Message, error)
	SetMessage(message *Message) error
}
