package config

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/postgres"
	"github.com/wtg/shuttletracker/smooth"
	"github.com/wtg/shuttletracker/spoofer"
	"github.com/wtg/shuttletracker/updater"
)

// Config is the global configuration struct.
type Config struct {
	Updater               *updater.Config
	API                   *api.Config
	Log                   *log.Config
	Postgres              *postgres.Config
	SmoothTrackingManager *smooth.Config
	Spoofer               *spoofer.Config
}

// New creates a new, global Config. Reads in configuration from config files.
func New() (*Config, error) {
	cfg := &Config{}

	// Create a global viper. Eventually, we should be creating sub-vipers and passing them into each NewConfig(),
	// but I have been unsuccessful in getting the resulting sub-vipers merged into one viper.
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg.API = api.NewConfig(v)
	cfg.Updater = updater.NewConfig(v)
	cfg.Log = log.NewConfig(v)
	cfg.SmoothTrackingManager = smooth.NewConfig(v)
	cfg.Spoofer = spoofer.NewConfig(v)

	pgCfg, err := postgres.NewConfig(v)
	if err != nil {
		return nil, err
	}
	cfg.Postgres = pgCfg

	v.SetConfigName("conf")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		log.Info("No config file found; only reading from environment")
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// I have no idea why, but these two configs need to be reset after reading the file
	cfg.SmoothTrackingManager = smooth.BackupConfig(v)
	cfg.Spoofer = spoofer.BackupConfig(v)

	// Special case for setting log level after reading config
	log.SetLevel(cfg.Log.Level)
	log.Debugf("All settings: %+v", v.AllSettings())
	log.Debugf("API configuration: %+v", cfg.API)
	log.Debugf("Updater configuration: %+v", cfg.Updater)
	log.Debugf("Log configuration: %+v", cfg.Log)
	log.Debugf("Postgres configuration: %+v", cfg.Postgres)
	log.Debugf("Smooth tracking configuration: %+v", cfg.SmoothTrackingManager)
	log.Debugf("Spoofer configuration: %+v", cfg.Spoofer)

	return cfg, nil
}
