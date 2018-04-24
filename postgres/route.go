package postgres

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/lib/pq"
	"github.com/wtg/shuttletracker"
)

// pointsRegex matches points in a Postgres path, e.g. [(42.72283,-73.67964),(42.72297,-73.67948)]
var pointsRegex = regexp.MustCompile(`\((-?\d+\.?\d*),(-?\d+\.?\d*)\)`)

// RouteService implements shuttletracker.RouteService.
type RouteService struct {
	db *sql.DB
}

func (rs *RouteService) initializeSchema(db *sql.DB) error {
	rs.db = db
	schema := `
CREATE TABLE IF NOT EXISTS routes (
    id serial PRIMARY KEY,
	name text NOT NULL,
	created timestamp with time zone NOT NULL DEFAULT now(),
	updated timestamp with time zone NOT NULL DEFAULT now(),
	enabled boolean NOT NULL,
	width smallint NOT NULL DEFAULT 4,
	color varchar(9) NOT NULL DEFAULT '#ffffff',
	points path
);
CREATE TABLE IF NOT EXISTS routes_stops (
	id serial PRIMARY KEY,
	route_id integer REFERENCES routes ON DELETE CASCADE NOT NULL,
	stop_id integer REFERENCES stops NOT NULL,
	"order" integer NOT NULL,
	UNIQUE (route_id, "order")
);`
	_, err := rs.db.Exec(schema)
	return err
}

// TODO: document this
type scanPoints struct {
	points []shuttletracker.Point
}

// TODO: document this
func (p *scanPoints) Scan(src interface{}) error {
	if src == nil {
		p.points = []shuttletracker.Point{}
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return errors.New("unable to scan points")
	}

	for _, pointMatch := range pointsRegex.FindAllSubmatch(b, -1) {
		lat, err := strconv.ParseFloat(string(pointMatch[1]), 64)
		if err != nil {
			return err
		}
		lon, err := strconv.ParseFloat(string(pointMatch[2]), 64)
		if err != nil {
			return err
		}
		point := shuttletracker.Point{
			Latitude:  lat,
			Longitude: lon,
		}
		p.points = append(p.points, point)
	}
	return nil
}

// Routes returns all Routes in the database.
func (rs *RouteService) Routes() ([]*shuttletracker.Route, error) {
	routes := []*shuttletracker.Route{}
	query := "SELECT r.id, r.name, r.created, r.updated, r.enabled, r.width, r.color, r.points," +
		" array_remove(array_agg(rs.stop_id ORDER BY rs.order ASC), NULL) as stop_ids" +
		" FROM routes r LEFT JOIN routes_stops rs" +
		" ON r.id = rs.route_id GROUP BY r.id;"
	rows, err := rs.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		r := &shuttletracker.Route{}
		p := scanPoints{}
		err := rows.Scan(&r.ID, &r.Name, &r.Created, &r.Updated, &r.Enabled, &r.Width, &r.Color, &p, pq.Array(&r.StopIDs))
		if err != nil {
			return nil, err
		}
		r.Points = p.points
		routes = append(routes, r)
	}
	return routes, nil
}

// Route returns the Route with the provided ID.
func (rs *RouteService) Route(id int) (*shuttletracker.Route, error) {
	query := "SELECT r.id, r.name, r.created, r.updated, r.enabled, r.width, r.color, r.points," +
		" array_remove(array_agg(rs.stop_id ORDER BY rs.order ASC), NULL) as stop_ids" +
		" FROM routes r LEFT JOIN routes_stops rs" +
		" ON r.id = rs.route_id WHERE r.id = $1 GROUP BY r.id;"
	row := rs.db.QueryRow(query, id)
	r := &shuttletracker.Route{}
	p := scanPoints{}
	err := row.Scan(&r.Name, &r.Created, &r.Updated, &r.Enabled, &r.Width, &r.Color, &p, pq.Array(&r.StopIDs))
	if err != nil {
		return nil, err
	}
	r.Points = p.points
	return r, nil
}

// TODO: document this
type valuePoints []shuttletracker.Point

// TODO: document this
func (p valuePoints) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, nil
	}

	buf := &bytes.Buffer{}
	err := buf.WriteByte('[')
	if err != nil {
		return nil, err
	}

	for i, point := range p {
		_, err = buf.WriteString(fmt.Sprintf("(%f,%f)", point.Latitude, point.Longitude))
		if err != nil {
			return nil, err
		}
		if i != len(p)-1 {
			err = buf.WriteByte(',')
			if err != nil {
				return nil, err
			}
		}
	}
	err = buf.WriteByte(']')
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// CreateRoute creates a Route.
func (rs *RouteService) CreateRoute(route *shuttletracker.Route) error {
	statement := "INSERT INTO routes (name, enabled, width, color, points)" +
		" VALUES ($1, $2, $3, $4, $5) RETURNING id, created, updated;"
	row := rs.db.QueryRow(statement, route.Name, route.Enabled, route.Width, route.Color, valuePoints(route.Points))
	return row.Scan(&route.ID, &route.Created, &route.Updated)
}

// DeleteRoute deletes a Route.
func (rs *RouteService) DeleteRoute(id int) error {
	statement := "DELETE FROM routes WHERE id = $1;"
	_, err := rs.db.Exec(statement, id)
	return err
}

// ModifyRoute modifies an existing Route.
func (rs *RouteService) ModifyRoute(route *shuttletracker.Route) error {
	statement := "UPDATE routes SET name = $1, enabled = $2, width = $3, color = $4, points = $5, updated = now()" +
		" WHERE id = $6 RETURNING updated;"
	row := rs.db.QueryRow(statement, route.Name, route.Enabled, route.Width, route.Color, valuePoints(route.Points), route.ID)
	return row.Scan(&route.Updated)
}
