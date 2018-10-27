package postgres

import (
	"database/sql"

	"github.com/spf13/viper"
)

/*
Postgres implements shuttletracker.VehicleService, shuttletracker.RouteService,
shuttletracker.StopService, shuttletracker.LoctionService, shuttletracker.MessageService,
and shuttletracker.UserService.
*/
type Postgres struct {
	VehicleService
	RouteService
	StopService
	LocationService
	MessageService
	UserService
}

// Config contains database connection information.
type Config struct {
	URL string
}

// New returns a configured Postgres.
func New(cfg Config) (*Postgres, error) {
	db, err := sql.Open("postgres", cfg.URL)
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
	err = pg.StopService.initializeSchema(db)
	if err != nil {
		return nil, err
	}
	err = pg.RouteService.initializeSchema(db)
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

// NewConfig creates a new Config.
func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		URL: "postgres://localhost/shuttletracker?sslmode=disable",
	}
	v.SetDefault("postgres.url", cfg.URL)

	// Allow DATABASE_URL to set the Postgres connection string for ease of deployment.
	v.BindEnv("postgres.url", "DATABASE_URL")

	return cfg
}
