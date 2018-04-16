// Package main bundles together all of shuttletracker's subpackages
// to create, configure, and run the shuttle tracker.
package main

import (
	"github.com/kochman/runner"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/config"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/postgres"
	"github.com/wtg/shuttletracker/updater"
)

func main() {
	Run()
}

// Run starts the Shuttle Tracker and blocks forever.
func Run() {
	log.Info("Shuttle Tracker starting...")

	// Config
	cfg, err := config.New()
	if err != nil {
		log.WithError(err).Error("Could not create config.")
		return
	}

	runner := runner.New()

	// Log
	log.SetLevel(cfg.Log.Level)

	// Database
	db, err := database.NewMongoDB(*cfg.Database)
	if err != nil {
		log.WithError(err).Errorf("MongoDB connection to \"%v\" failed.", cfg.Database.MongoURL)
		return
	}

	// Model service
	var ms shuttletracker.ModelService
	ms, err = postgres.NewModelService("postgres://localhost/shuttletracker?sslmode=disable")
	if err != nil {
		log.WithError(err).Error("unable to create ModelService")
		return
	}

	// Make shuttle position updater
	updater, err := updater.New(*cfg.Updater, ms)
	if err != nil {
		log.WithError(err).Error("Could not create updater.")
		return
	}
	runner.Add(updater)

	// Make API server
	api, err := api.New(*cfg.API, db, ms)
	if err != nil {
		log.WithError(err).Error("Could not create API server.")
		return
	}
	runner.Add(api)

	// Run all runnables
	runner.Run()
}
