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

	// Start shuttle position updater
	updater := updater.New(*cfg.Updater, *db)
	go updater.Run()

	// Start API server
	api := api.New(*cfg.API, *db)
	go api.Run()

	// Wait until quit
	quit := make(chan bool, 0)
	select {
	case <-quit:
	}
}
