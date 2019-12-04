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
	schema := `
CREATE TABLE IF NOT EXISTS messages (
	id bool PRIMARY KEY DEFAULT true CHECK (id = true),
	type text,
	message text,
	enabled bool NOT NULL,
	created timestamp with time zone NOT NULL DEFAULT now(),
	updated timestamp with time zone NOT NULL DEFAULT now(),
	link text
);
ALTER TABLE messages ADD COLUMN IF NOT EXISTS link text;`
	_, err := ms.db.Exec(schema)
	return err
}

// Message returns the Message.
func (ms *MessageService) Message() (*shuttletracker.Message, error) {
	query := "SELECT type, message, enabled, created, updated, link FROM messages;"
	row := ms.db.QueryRow(query)
	message := &shuttletracker.Message{}
	err := row.Scan(&message.Message, &message.Enabled, &message.Created, &message.Updated, &message.Link)
	if err == sql.ErrNoRows {
		return nil, shuttletracker.ErrMessageNotFound
	} else if err != nil {
		return nil, err
	}
	return message, nil
}

// func (ms *MessageService) AnnouncementMessage() (*shuttleTracker.Message, error) {
// 	query := "SELECT "
// }

// SetMessage updates the Message.
func (ms *MessageService) SetMessage(message *shuttletracker.Message) error {
	statement := "INSERT INTO messages (type, message, enabled, updated, link) VALUES ($1, $2, now(), $3)" +
		" ON CONFLICT (id) DO UPDATE SET type=excluded.type, message = excluded.message, enabled = excluded.enabled, updated = excluded.updated, link = excluded.link" +
		" RETURNING created, updated;"
	row := ms.db.QueryRow(statement, message.Message, message.Enabled, message.Link)
	return row.Scan(&message.Created, &message.Updated)
}
