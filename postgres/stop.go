package postgres

import (
	"database/sql"

	"github.com/wtg/shuttletracker"
)

// StopService is an implementation of shuttletracker.StopService.
type StopService struct {
	db *sql.DB
}

func (ss *StopService) initializeSchema(db *sql.DB) error {
	ss.db = db
	schema := ``
	_, err := ss.db.Exec(schema)
	return err
}

func (ss *StopService) CreateStop(stop *shuttletracker.Stop) error {
	return nil
}

func (ss *StopService) Stops() ([]*shuttletracker.Stop, error) {
	return nil, nil
}

func (ss *StopService) DeleteStop(id int) error {
	return nil
}
