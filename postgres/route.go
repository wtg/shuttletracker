package postgres

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/lib/pq"

	"github.com/wtg/shuttletracker"
)

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
	color varchar(9) NOT NULL DEFAULT '#ffffff'
);
CREATE TABLE IF NOT EXISTS routes_points (
	id serial PRIMARY KEY,
	latitude double precision NOT NULL,
	longitude double precision NOT NULL,
	route_id integer REFERENCES routes NOT NULL
);`
	_, err := rs.db.Exec(schema)
	return err
}

type scanPoint struct {
	latitude, longitude float64
}

type scanPoints []scanPoint

func (p scanPoints) points() []shuttletracker.Point {
	points := []shuttletracker.Point{}
	for _, point := range p {
		newPoint := shuttletracker.Point{
			Latitude:  point.latitude,
			Longitude: point.longitude,
		}
		points = append(points, newPoint)
	}
	return points
}

func (p *scanPoint) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("unable to scan point")
	}

	s := string(b)
	s = s[1 : len(s)-1]
	split := strings.Split(s, ",")

	lat, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		return err
	}
	p.latitude = lat

	lon, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return err
	}
	p.longitude = lon

	return nil
}

func (rs *RouteService) Routes() ([]*shuttletracker.Route, error) {
	routes := []*shuttletracker.Route{}
	query := "SELECT r.id, r.name, r.created, r.updated, r.enabled, r.width, r.color, array_agg(point(rp.latitude, rp.longitude)) as points" +
		" FROM routes r" +
		" LEFT JOIN routes_points rp ON r.id = rp.route_id" +
		" GROUP BY r.id;"
	rows, err := rs.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		r := &shuttletracker.Route{}
		p := scanPoints{}
		err := rows.Scan(&r.ID, &r.Name, &r.Created, &r.Updated, &r.Enabled, &r.Width, &r.Color, pq.Array(&p))
		if err != nil {
			return nil, err
		}
		r.Points = p.points()
		routes = append(routes, r)
	}
	return routes, nil
}

func (rs *RouteService) Route(id int) (*shuttletracker.Route, error) {
	return nil, nil
}

func (rs *RouteService) CreateRoute(route *shuttletracker.Route) error {
	return nil
}

func (rs *RouteService) DeleteRoute(id int) error {
	return nil
}

func (rs *RouteService) ModifyRoute(route *shuttletracker.Route) error {
	return nil
}
