package postgres

import (
	"database/sql"

	// Postgres driver for database/sql
	_ "github.com/lib/pq"

	"github.com/wtg/shuttletracker"
)

// RouteService implements shuttletracker.RouteService.
type RouteService struct {
	db *sql.DB
}

func (r *RouteService) initializeSchema(db *sql.DB) error {
	r.db = db
	schema := `
--DROP TABLE routes;
CREATE TABLE IF NOT EXISTS routes (
    id serial PRIMARY KEY,
	name text,
	created timestamp with time zone NOT NULL DEFAULT now(),
	updated timestamp with time zone NOT NULL DEFAULT now(),
	enabled boolean NOT NULL,
	width integer NOT NULL,
	color string
);
DROP TABLE routes_points;
CREATE TABLE IF NOT EXISTS routes_points (
	id serial PRIMARY KEY,
	latitude double precision NOT NULL,
	longitude double precision NOT NULL,
	route integer REFERENCES routes NOT NULL
);
    `
	_, err := r.db.Exec(schema)
	return err
}

func (rs *RouteService) Routes() ([]*shuttletracker.Route, error) {
	return nil, nil
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
