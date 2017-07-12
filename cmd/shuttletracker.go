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

	runner := NewRunner()

	// Log
	log.SetLevel(cfg.Log.Level)

	// Database
	db, err := database.New(*cfg.Database)
	if err != nil {
		log.WithError(err).Errorf("MongoDB connection to \"%v\" failed.", cfg.Database.MongoURL)
		return
	}

	// Make shuttle position updater
	updater, err := updater.New(*cfg.Updater, *db)
	if err != nil {
		log.WithError(err).Error("Could not create updater.")
		return
	}
	runner.Add(updater)

	// Make API server
	api, err := api.New(*cfg.API, *db)
	if err != nil {
		log.WithError(err).Error("Could not create API server.")
		return
	}
	runner.Add(api)

	// Run all runnables
	runner.Run()
}
