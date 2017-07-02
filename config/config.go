package config

import (
	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/updater"

	"github.com/fatih/structs"
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
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
	}

	// Dynamically determine what configuration options we've got.
	cfgStruct := structs.New(&cfg)
	for pkg, str := range cfgStruct.Map() {
		if structs.IsStruct(str) {
			str := structs.New(str)
			for _, field := range str.Fields() {
				option := field.Name()
				viperPath := pkg + "." + option
				if viper.IsSet(viperPath) {
					field.Set(viper.Get(viperPath))
				}
			}
		}
	}

	return cfg
}
