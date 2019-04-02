package postgres

import (
	"database/sql"
	"time"

	"github.com/wtg/shuttletracker"
)

// LocationService implements shuttletracker.LocationService.
type LocationService struct {
	db *sql.DB
}

func (ls *LocationService) initializeSchema(db *sql.DB) error {
	ls.db = db
	schema := `
CREATE TABLE IF NOT EXISTS locations (
	id serial PRIMARY KEY,
	tracker_id varchar(10) NOT NULL,
	latitude double precision NOT NULL,
	longitude double precision NOT NULL,
	heading real NOT NULL,
	speed real NOT NULL,
	time timestamp with time zone NOT NULL,
	route_id integer,
	created timestamp with time zone NOT NULL DEFAULT now(),
	UNIQUE (tracker_id, time)
);`
	_, err := ls.db.Exec(schema)
	return err
}

// CreateLocation creates a Location in the database.
func (ls *LocationService) CreateLocation(l *shuttletracker.Location) error {
	query := `
WITH location AS (
	INSERT INTO locations (
		tracker_id,
		latitude,
		longitude,
		heading,
		speed,
		time,
		route_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, tracker_id, created)
SELECT
	location.id AS location_id,
	vehicles.id AS vehicle_id,
	location.created
FROM location
LEFT JOIN vehicles ON vehicles.tracker_id = location.tracker_id;`
	row := ls.db.QueryRow(query, l.TrackerID, l.Latitude, l.Longitude, l.Heading, l.Speed, l.Time, l.RouteID)
	err := row.Scan(&l.ID, &l.VehicleID, &l.Created)
	return err
}

// DeleteLocationsBefore deletes all Locations in the database with tracker times before the provided Time.
func (ls *LocationService) DeleteLocationsBefore(before time.Time) (int, error) {
	statement := "DELETE FROM locations WHERE time < $1;"
	res, err := ls.db.Exec(statement, before)
	if err != nil {
		return 0, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

// LocationsSince returns all Locations since a tracker Time for a certain Vehicle, ordered newest to oldest.
func (ls *LocationService) LocationsSince(vehicleID int64, since time.Time) ([]*shuttletracker.Location, error) {
	locations := []*shuttletracker.Location{}
	query := "SELECT l.id, l.tracker_id, l.latitude, l.longitude, l.heading, l.speed, l.time, l.route_id, l.created " +
		"FROM locations l, vehicles v WHERE l.tracker_id = v.tracker_id AND v.id = $1 AND l.time > $2 ORDER BY l.created DESC;"
	rows, err := ls.db.Query(query, vehicleID, since)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		l := &shuttletracker.Location{
			VehicleID: &vehicleID,
		}
		err := rows.Scan(&l.ID, &l.TrackerID, &l.Latitude, &l.Longitude, &l.Heading, &l.Speed, &l.Time, &l.RouteID, &l.Created)
		if err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}

// LatestLocation returns the most recent Location created for a Vehicle.
func (ls *LocationService) LatestLocation(vehicleID int64) (*shuttletracker.Location, error) {
	l := &shuttletracker.Location{
		VehicleID: &vehicleID,
	}
	query := "SELECT l.id, l.tracker_id, l.latitude, l.longitude, l.heading, l.speed, l.time, l.route_id, l.created " +
		"FROM locations l, vehicles v WHERE l.tracker_id = v.tracker_id AND v.id = $1 " +
		"ORDER BY l.created DESC LIMIT 1;"
	row := ls.db.QueryRow(query, vehicleID)
	err := row.Scan(&l.ID, &l.TrackerID, &l.Latitude, &l.Longitude, &l.Heading, &l.Speed, &l.Time, &l.RouteID, &l.Created)
	if err == sql.ErrNoRows {
		return nil, shuttletracker.ErrLocationNotFound
	} else if err != nil {
		return nil, err
	}
	return l, nil
}
