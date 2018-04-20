package postgres

import "database/sql"

/*
Postgres implements shuttletracker.VehicleService, shuttletracker.RouteService,
shuttletracker.LoctionService, and shuttletracker.MessageService.
*/
type Postgres struct {
	db *sql.DB
	VehicleService
	RouteService
	StopService
	LocationService
	MessageService
	UserService
}

// New returns a configured Postgres.
func New(url string) (*Postgres, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	pg := &Postgres{
		db: db,
	}

	err = pg.VehicleService.initializeSchema()
	if err != nil {
		return nil, err
	}
	err = pg.RouteService.initializeSchema()
	if err != nil {
		return nil, err
	}
	err = pg.StopService.initializeSchema()
	if err != nil {
		return nil, err
	}
	err = pg.LocationService.initializeSchema()
	if err != nil {
		return nil, err
	}
	err = pg.MessageService.initializeSchema()
	if err != nil {
		return nil, err
	}
	err = pg.UserService.initializeSchema()
	if err != nil {
		return nil, err
	}

	return pg, nil
}
