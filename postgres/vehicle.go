package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/wtg/shuttletracker"
)

type VehicleService struct {
	db *sql.DB
}

// NewVehicleService returns a configured VehicleService.
func NewVehicleService(url string) (*VehicleService, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	vs := &VehicleService{
		db: db,
	}

	err = vs.initializeSchema()
	if err != nil {
		return nil, err
	}

	return vs, nil
}

func (v *VehicleService) initializeSchema() error {
	schema := `
CREATE TABLE vehicles (
    id serial primary key,
    name text
);
    `
	_, err := v.db.Exec(schema)
	return err
}

// CreateVehicle creates a Vehicle.
func (v *VehicleService) CreateVehicle(vehicle *shuttletracker.Vehicle) error {
	return nil
}

// DeleteVehicle deletes a Vehicle by its ID.
func (v *VehicleService) DeleteVehicle(id int) error {
	return nil
}

// GetVehicle returns a Vehicle by its ID.
func (v *VehicleService) Vehicle(id int) (*shuttletracker.Vehicle, error) {
	var vehicle *shuttletracker.Vehicle
	return vehicle, nil
}

// GetVehicles returns all Vehicles.
func (v *VehicleService) Vehicles() ([]*shuttletracker.Vehicle, error) {
	var vehicles []*shuttletracker.Vehicle
	return vehicles, nil
}

// GetVehicle returns a Vehicle by its ID.
func (v *VehicleService) EnabledVehicles() ([]*shuttletracker.Vehicle, error) {
	var vehicles []*shuttletracker.Vehicle
	return vehicles, nil
}

// GetEnabledVehicles returns all Vehicles that are enabled.
func (v *VehicleService) GetEnabledVehicles() ([]*shuttletracker.Vehicle, error) {
	var vehicles []*shuttletracker.Vehicle
	return vehicles, nil
}

// ModifyVehicle updates a Vehicle by its ID.
func (v *VehicleService) ModifyVehicle(vehicle *shuttletracker.Vehicle) error {
	return nil
}

func (v *VehicleService) VehicleWithTrackerID(id int) (*shuttletracker.Vehicle, error) {
	var vehicle *shuttletracker.Vehicle
	return vehicle, nil
}
