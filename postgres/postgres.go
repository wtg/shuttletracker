package postgres

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
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

	listener := pq.NewListener(cfg.URL, time.Second, time.Minute, nil)

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
	err = pg.LocationService.initializeSchema(db, listener)
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

	go pg.LocationService.run()

	return pg, nil
}

// NewConfig creates a new Config.
func NewConfig(v *viper.Viper) (*Config, error) {
	cfg := &Config{
		URL: "postgres://localhost/shuttletracker?sslmode=disable",
	}
	v.SetDefault("postgres.url", cfg.URL)

	// Allow DATABASE_URL to set the Postgres connection string for ease of deployment.
	err := v.BindEnv("postgres.url", "DATABASE_URL")
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
