package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/config"
	"github.com/wtg/shuttletracker/postgres"
)

var Add bool

func init() {
	adminsCmd.Flags().BoolVar(&Add, "add", false, "add administrator")

	rootCmd.AddCommand(adminsCmd)
}

var adminsCmd = &cobra.Command{
	Use:   "admins",
	Short: "Manage Shuttle Tracker administrators",
	Args: func(cms *cobra.Command, args []string) error {
		if Add && len(args) != 1 {
			return errors.New("expects exactly one argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Config
		cfg, err := config.New()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to read configuration.")
			os.Exit(1)
		}

		pg, err := postgres.New(*cfg.Postgres)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to connect to Postgres:", err)
			os.Exit(1)
		}
		var us shuttletracker.UserService = pg

		if Add {
			username := args[0]
			user := &shuttletracker.User{
				Username: username,
			}
			err := us.CreateUser(user)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Unable to add admin:", err)
				os.Exit(1)
			}
			fmt.Printf("Added %s.\n", username)
		} else {
			users, err := us.Users()
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to get users:", err)
				os.Exit(1)
			}

			if len(users) == 0 {
				fmt.Println("No Shuttle Tracker administrators.")
				return
			}

			for _, user := range users {
				fmt.Println(user.Username)
			}
		}
	},
}
