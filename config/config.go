package config

import (
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/updater"
	"github.com/wtg/shuttletracker/api"
)

type Config struct {
	DataFeed             string `env:"DATA_FEED"`
	UpdateInterval       int    `env:"UPDATE_INTERVAL" envDefault:"15"`
	MongoURL             string `env:"MONGO_URL" envDefault:"localhost:27017"`
	GoogleMapAPIKey      string
	GoogleMapMinDistance int
	CasURL               string `env:"CAS_URL"`
	Authenticate         bool   `env:"AUTHENTICATE" envDefault:"true"`
	Database database.Config
	Updater updater.Config
	API api.Config
}

func New() *Config {
	return &Config{
		Database: database.Config{
			MongoURL: "localhost",
		},
		Updater: updater.Config{

		},
		API: api.Config{},
	}
}