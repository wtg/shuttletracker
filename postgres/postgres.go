package postgres

import "database/sql"

// ModelService implements shuttletracker.VehicleService.
type ModelService struct {
	VehicleService
	RouteService
	LocationService
	db *sql.DB
}

// NewModelService returns a configured ModelService.
func NewModelService(url string) (*ModelService, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	ms := &ModelService{
		db: db,
	}

	err = ms.VehicleService.initializeSchema()
	if err != nil {
		return nil, err
	}
	err = ms.RouteService.initializeSchema()
	if err != nil {
		return nil, err
	}

	return ms, nil
}
