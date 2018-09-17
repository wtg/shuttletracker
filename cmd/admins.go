package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/config"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/postgres"
)

func init() {
	rootCmd.AddCommand(adminsCmd)
}

var adminsCmd = &cobra.Command{
	Use:   "admins",
	Short: "Manage Shuttle Tracker administrators",
	Run: func(cmd *cobra.Command, args []string) {
		// Config
		cfg, err := config.New()
		if err != nil {
			log.WithError(err).Error("unable to read configuration")
			os.Exit(1)
		}

		pg, err := postgres.New(*cfg.Postgres)
		if err != nil {
			log.WithError(err).Error("unable to create Postgres")
			os.Exit(1)
		}

		var us shuttletracker.UserService = pg
		users, err := us.Users()
		if err != nil {
			log.WithError(err).Error("unable to get users")
			os.Exit(1)
		}

		if len(users) == 0 {
			fmt.Println("No Shuttle Tracker administrators.")
			return
		}

		for _, user := range users {
			fmt.Println(user.Username)
		}
	},
}
