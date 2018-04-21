package postgres

import "database/sql"

/*
Postgres implements shuttletracker.VehicleService, shuttletracker.RouteService,
shuttletracker.LoctionService, and shuttletracker.MessageService.
*/
type Postgres struct {
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

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	pg := &Postgres{}

	err = pg.VehicleService.initializeSchema(db)
	if err != nil {
		return nil, err
	}
	err = pg.RouteService.initializeSchema(db)
	if err != nil {
		return nil, err
	}
	err = pg.StopService.initializeSchema(db)
	if err != nil {
		return nil, err
	}
	err = pg.LocationService.initializeSchema(db)
	if err != nil {
		return nil, err
	}
	err = pg.MessageService.initializeSchema(db)
	if err != nil {
		return nil, err
	}
	err = pg.UserService.initializeSchema(db)
	if err != nil {
		return nil, err
	}

	return pg, nil
}
