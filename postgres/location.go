package postgres

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/lib/pq"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
)

const locationsInsertChannel = "locations.insert"

// LocationService implements shuttletracker.LocationService.
type LocationService struct {
	db          *sql.DB
	listener    *pq.Listener
	addSub      chan chan *shuttletracker.Location
	subscribers []chan *shuttletracker.Location
}

func (ls *LocationService) initializeSchema(db *sql.DB, listener *pq.Listener) error {
	ls.db = db
	ls.listener = listener
	ls.addSub = make(chan chan *shuttletracker.Location)
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
);

-- notify clients when locations inserted
CREATE OR REPLACE FUNCTION locations_insert_notify() RETURNS trigger AS $$
BEGIN
        PERFORM pg_notify('locations.insert', NEW.id::text);
        RETURN NEW;
END
$$ LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS locations_insert on locations;
CREATE TRIGGER locations_insert AFTER INSERT ON locations FOR EACH ROW EXECUTE PROCEDURE locations_insert_notify();
`
	_, err := ls.db.Exec(schema)
	return err
}

func (ls *LocationService) run() {
	err := ls.listener.Listen(locationsInsertChannel)
	if err != nil {
		log.WithError(err).Error("unable to listen for inserted locations")
		return
	}

	for {
		select {
		case c := <-ls.addSub:
			ls.subscribers = append(ls.subscribers, c)
		case n := <-ls.listener.Notify:
			if n == nil {
				break
			}
			if n.Channel != locationsInsertChannel {
				continue
			}

			id, err := strconv.ParseInt(n.Extra, 10, 64)
			if err != nil {
				log.WithError(err).Error("unable to parse location ID")
				continue
			}

			loc, err := ls.Location(id)
			if err != nil {
				log.WithError(err).Error("unable to get location")
				continue
			}
			for _, sub := range ls.subscribers {
				sub <- loc
			}
		}
	}
}

// SubscribeLocations returns a chan that receives each new Location after it is
// written to the database.
func (ls *LocationService) SubscribeLocations() chan *shuttletracker.Location {
	c := make(chan *shuttletracker.Location)
	ls.addSub <- c
	return c
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

// LatestLocations returns the most recent Location created for all enabled Vehicles.
func (ls *LocationService) LatestLocations() ([]*shuttletracker.Location, error) {
	locations := []*shuttletracker.Location{}
	query := `
SELECT l.id, l.tracker_id, l.latitude, l.longitude, l.heading, l.speed, l.time, l.route_id, l.created, v.id
FROM vehicles v,
        locations l JOIN (
                SELECT tracker_id, max(created) AS created
                from locations
                group by tracker_id) AS l2
        ON l.tracker_id = l2.tracker_id AND l.created = l2.created
WHERE v.tracker_id = l.tracker_id;
	`
	rows, err := ls.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		l := &shuttletracker.Location{}
		err := rows.Scan(&l.ID, &l.TrackerID, &l.Latitude, &l.Longitude, &l.Heading, &l.Speed, &l.Time, &l.RouteID, &l.Created, &l.VehicleID)
		if err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}

// Location returns a Location with the provided ID.
func (ls *LocationService) Location(id int64) (*shuttletracker.Location, error) {
	l := &shuttletracker.Location{
		ID: id,
	}
	query := "SELECT l.tracker_id, l.latitude, l.longitude, l.heading, l.speed, l.time, l.route_id, l.created, v.id " +
		"FROM locations l, vehicles v WHERE l.tracker_id = v.tracker_id AND l.id = $1;"
	row := ls.db.QueryRow(query, id)
	err := row.Scan(&l.TrackerID, &l.Latitude, &l.Longitude, &l.Heading, &l.Speed, &l.Time, &l.RouteID, &l.Created, &l.VehicleID)
	if err == sql.ErrNoRows {
		return nil, shuttletracker.ErrLocationNotFound
	} else if err != nil {
		return nil, err
	}
	return l, nil
}
