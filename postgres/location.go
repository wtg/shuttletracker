package postgres

import (
	"database/sql"
	"time"

	"github.com/wtg/shuttletracker"
)

type LocationService struct {
	db *sql.DB
}

func (ls *LocationService) initializeSchema(db *sql.DB) error {
	ls.db = db
	schema := `
--DROP TABLE locations;
CREATE TABLE IF NOT EXISTS locations (
	id serial PRIMARY KEY,
	tracker_id varchar(10) NOT NULL,
	latitude double precision NOT NULL,
	longitude double precision NOT NULL,
	heading real NOT NULL,
	speed real NOT NULL,
	time timestamp with time zone NOT NULL,
	route_id integer,
	created timestamp with time zone NOT NULL DEFAULT now()
);
	`
	_, err := ls.db.Exec(schema)
	return err
}

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
	location.id as location_id,
	vehicles.id AS vehicle_id,
	location.created
FROM location
LEFT JOIN vehicles ON vehicles.tracker_id = location.tracker_id;`
	row := ls.db.QueryRow(query, l.TrackerID, l.Latitude, l.Longitude, l.Heading, l.Speed, l.Time, l.RouteID)
	err := row.Scan(&l.ID, &l.VehicleID, &l.Created)
	return err
}

func (ls *LocationService) DeleteLocationsBefore(before time.Time) (int, error) {
	return 0, nil
}

func (ls *LocationService) LocationsSince(vehicleID int, since time.Time) ([]*shuttletracker.Location, error) {
	return nil, nil
}

func (ls *LocationService) LatestLocation(vehicleID int) (*shuttletracker.Location, error) {
	l := &shuttletracker.Location{
		VehicleID: &vehicleID,
	}
	query := "SELECT l.id, l.tracker_id, l.latitude, l.longitude, l.heading, l.speed, l.time, l.route_id " +
		"FROM locations l, vehicles v WHERE l.tracker_id = v.tracker_id AND v.id = $1 " +
		"ORDER BY l.created DESC;"
	row := ls.db.QueryRow(query, vehicleID)
	err := row.Scan(&l.ID, &l.TrackerID, &l.Latitude, &l.Longitude, &l.Heading, &l.Speed, &l.Time, &l.RouteID)
	if err == sql.ErrNoRows {
		return nil, shuttletracker.ErrLocationNotFound
	} else if err != nil {
		return nil, err
	}
	return l, nil
}
