package config

import (
	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/updater"

	"github.com/spf13/viper"
)

// Global configuration struct.
type Config struct {
	Database *database.Config `structs:",omitnested"`
	Updater  *updater.Config  `structs:",omitnested"`
	API      *api.Config      `structs:",omitnested"`
	Log      *log.Config      `structs:",omitnested"`
}

// Create a new, global Config. Reads in configuration from the environment and config files.
func New() Config {
	cfg := Config{
		Database: database.NewConfig(),
		Updater:  updater.NewConfig(),
		API:      api.NewConfig(),
		Log:      log.NewConfig(),
	}

	viper.SetConfigName("conf")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Error("Unable to read configuration.")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.WithError(err).Error("Unable to unmarshal configuration.")
	}

	return cfg
}
