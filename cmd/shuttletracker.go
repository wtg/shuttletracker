// Package cmd bundles together all of shuttletracker's subpackages
// to create, configure, and run the shuttle tracker.
package cmd

import (
	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/config"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/updater"
)

// Run starts the shuttle tracker and blocks forever.
func Run() {
	log.Info("Shuttle Tracker starting...")

	// Config
	cfg, err := config.New()
	if err != nil {
		log.WithError(err).Error("Could not create config.")
		return
	}

	// Log
	log.SetLevel(cfg.Log.Level)

	// Database
	db := database.New(*cfg.Database)

	// Start shuttle position updater
	updater, err := updater.New(*cfg.Updater, *db)
	if err != nil {
		log.WithError(err).Error("Could not create updater.")
		return
	}
	go updater.Run()

	// Start API server
	api, err := api.New(*cfg.API, *db)
	if err != nil {
		log.WithError(err).Error("Could not create API.")
		return
	}
	go api.Run()

	log.Info("Shuttle Tracker started.")

	// Wait until quit
	quit := make(chan bool, 0)
	select {
	case <-quit:
	}
}
