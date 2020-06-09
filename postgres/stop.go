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
	schema := `
CREATE TABLE IF NOT EXISTS stops (
	id serial PRIMARY KEY,
	name text,
	description text,
	latitude double precision NOT NULL,
	longitude double precision NOT NULL,
	created timestamp with time zone NOT NULL DEFAULT now(),
	updated timestamp with time zone NOT NULL DEFAULT now()
);`
	_, err := ss.db.Exec(schema)
	return err
}

// CreateStop creates a Stop.
func (ss *StopService) CreateStop(stop *shuttletracker.Stop) error {
	statement := "INSERT INTO stops (name, description, latitude, longitude) VALUES" +
		" ($1, $2, $3, $4) RETURNING id, created, updated;"
	row := ss.db.QueryRow(statement, stop.Name, stop.Description, stop.Latitude, stop.Longitude)
	return row.Scan(&stop.ID, &stop.Created, &stop.Updated)
}

// Stop returns a Stop by its ID.
func (ss *StopService) Stop(id int64) (*shuttletracker.Stop, error) {
	s := &shuttletracker.Stop{
		ID: id,
	}

	statement := "SELECT s.name, s.created, s.updated, s.description, s.latitude, s.longitude" +
		" FROM stops s WHERE id = $1;"
	row := ss.db.QueryRow(statement, id)
	err := row.Scan(&s.Name, &s.Created, &s.Updated, &s.Description, &s.Latitude, &s.Longitude)
	if err == sql.ErrNoRows {
		return nil, shuttletracker.ErrStopNotFound
	}

	return s, err
}

// Stops returns all Stops.
func (ss *StopService) Stops() ([]*shuttletracker.Stop, error) {
	stops := []*shuttletracker.Stop{}
	query := "SELECT s.id, s.name, s.created, s.updated, s.description, s.latitude, s.longitude" +
		" FROM stops s;"
	rows, err := ss.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &shuttletracker.Stop{}
		err := rows.Scan(&s.ID, &s.Name, &s.Created, &s.Updated, &s.Description, &s.Latitude, &s.Longitude)
		if err != nil {
			return nil, err
		}
		stops = append(stops, s)
	}
	return stops, nil
}

// DeleteStop deletes a Stop.
func (ss *StopService) DeleteStop(id int64) error {
	statement := "DELETE FROM stops WHERE id = $1;"
	result, err := ss.db.Exec(statement, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return shuttletracker.ErrStopNotFound
	}

	return nil
}
