package cmd

import (
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/updater"
	"github.com/wtg/shuttletracker/config"
	"github.com/wtg/shuttletracker/api"
)

func Run() {
	cfg := config.New()
	db := database.New(cfg.Database)
	updater := updater.New(cfg.Updater, *db)
	api := api.New(cfg.API, *db)

	// Start auto updater
	go updater.Run()
	api.Run()
}