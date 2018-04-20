package postgres

import (
	"database/sql"
	"time"

	"github.com/wtg/shuttletracker"
)

type LocationService struct {
	db *sql.DB
}

func (ls *LocationService) initializeSchema() error {
	schema := ``
	_, err := ls.db.Exec(schema)
	return err
}

func (ls *LocationService) CreateLocation(location *shuttletracker.Location) error {
	return nil
}

func (ls *LocationService) DeleteLocationsBefore(before time.Time) (int, error) {
	return 0, nil
}

func (ls *LocationService) LocationsSince(vehicleID int, since time.Time) ([]*shuttletracker.Location, error) {
	return nil, nil
}

func (ls *LocationService) LatestLocation(vehicleID int) (*shuttletracker.Location, error) {
	return nil, nil
}
