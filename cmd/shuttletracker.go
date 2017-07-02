package cmd

import (
	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/config"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/updater"
)

func Run() {
	log.Info("Shuttle Tracker starting...")
	cfg := config.New()

	log.SetLevel(cfg.Log.Level)

	db := database.New(*cfg.Database)
	log.Debug("Connected to database.")

	updater := updater.New(*cfg.Updater, *db)
	api := api.New(*cfg.API, *db)

	// Start auto updater
	go updater.Run()
	api.Run()
}
