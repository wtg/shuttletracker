package config

import (
	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/server"
	"github.com/wtg/shuttletracker/updater"

	"github.com/spf13/viper"
)

// Global configuration struct.
type Config struct {
	Database *database.Config
	Updater  *updater.Config
	API      *api.Config
	Log      *log.Config
	Server   *server.Config
}

// Create a new, global Config. Reads in configuration from config files.
func New() (*Config, error) {
	cfg := &Config{
		Database: database.NewConfig(),
		Updater:  updater.NewConfig(),
		API:      api.NewConfig(),
		Log:      log.NewConfig(),
		Server:   server.NewConfig(),
	}

	viper.SetConfigName("conf")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
