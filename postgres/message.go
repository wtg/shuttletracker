package postgres

import (
	"database/sql"

	"github.com/wtg/shuttletracker"
)

// MessageService implements a mock of shuttletracker.MessageService.
type MessageService struct {
	db *sql.DB
}

func (ms *MessageService) initializeSchema(db *sql.DB) error {
	ms.db = db
	schema := ``
	_, err := ms.db.Exec(schema)
	return err
}

// Message returns the Message.
func (ms *MessageService) Message() (*shuttletracker.Message, error) {
	return nil, nil
}

// SetMessage updates the Message.
func (ms *MessageService) SetMessage(message *shuttletracker.Message) error {
	return nil
}
