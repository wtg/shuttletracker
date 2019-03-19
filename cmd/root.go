package cmd

import (
	"fmt"
	"os"

	"github.com/kochman/runner"
	"github.com/spf13/cobra"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/config"
	"github.com/wtg/shuttletracker/eta"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/postgres"
	"github.com/wtg/shuttletracker/updater"
)

var rootCmd = &cobra.Command{
	Use:   "shuttletracker",
	Short: "Track RPI's shuttles",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Shuttle Tracker starting...")

		// Config
		cfg, err := config.New()
		if err != nil {
			log.WithError(err).Error("Could not create config.")
			return
		}

		runner := runner.New()

		pg, err := postgres.New(*cfg.Postgres)
		if err != nil {
			log.WithError(err).Error("unable to create Postgres")
			return
		}

		// Model service
		var ms shuttletracker.ModelService = pg

		// Message service
		var msg shuttletracker.MessageService = pg

		// User service
		var us shuttletracker.UserService = pg

		// Make shuttle position updater
		updater, err := updater.New(*cfg.Updater, ms)
		if err != nil {
			log.WithError(err).Error("Could not create updater.")
			return
		}
		runner.Add(updater)

		etaManager, err := eta.NewManager(ms, updater)
		if err != nil {
			log.WithError(err).Error("unable to create ETA manager")
			return
		}
		runner.Add(etaManager)

		// Make API server
		api, err := api.New(*cfg.API, ms, msg, us, updater, etaManager)
		if err != nil {
			log.WithError(err).Error("Could not create API server.")
			return
		}
		runner.Add(api)

		// Run all runnables
		runner.Run()
	},
}

// Execute makes the root command runnable.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
