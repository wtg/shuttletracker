// Package main bundles together all of shuttletracker's subpackages
// to create, configure, and run the shuttle tracker.
package main

import (
	"github.com/kochman/runner"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/config"
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

	pg, err := postgres.New(*cfg.Postgres)
	if err != nil {
		log.WithError(err).Error("unable to create Postgres")
		return
	}

	// Ensure service implementations actually implement their interfaces
	var ms shuttletracker.ModelService = pg
	var msg shuttletracker.MessageService = pg
	var us shuttletracker.UserService = pg

	// Make shuttle position updater
	updater, err := updater.New(*cfg.Updater, ms)
	if err != nil {
		log.WithError(err).Error("Could not create updater.")
		return
	}
	runner.Add(updater)

	// Make API server
	api, err := api.New(*cfg.API, ms, msg, us, updater)
	if err != nil {
		log.WithError(err).Error("Could not create API server.")
		return
	}
	runner.Add(api)

	// Run all runnables
	runner.Run()
}
