package config

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/postgres"
	"github.com/wtg/shuttletracker/updater"
)

// Config is the global configuration struct.
type Config struct {
	Updater  *updater.Config
	API      *api.Config
	Log      *log.Config
	Postgres *postgres.Config
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
	cfg.Log = log.NewConfig()
	cfg.Postgres = postgres.NewConfig(v)

	log.Debugf("All settings: %+v", v.AllSettings())

	v.SetConfigName("conf")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		log.Info("No config file found; only reading from environment")
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if v.IsSet("postgres.url") {
		log.Warn("POSTGRES_URL is deprecated and will be removed in the future. Use DATABASE_URL instead.")
	}

	return cfg, nil
}
