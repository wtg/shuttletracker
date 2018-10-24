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

// Add is a flag to put the admins command into "add" mode.
var Add bool

// Remove is a flag to put the admins command into "remove" mode.
var Remove bool

func init() {
	adminsCmd.Flags().BoolVar(&Add, "add", false, "add administrator")
	adminsCmd.Flags().BoolVar(&Remove, "remove", false, "remove administrator")

	rootCmd.AddCommand(adminsCmd)
}

var adminsCmd = &cobra.Command{
	Use:   "admins",
	Short: "Manage Shuttle Tracker administrators",
	Long:  "List, add, or remove Shuttle Tracker administrators by RCS ID.",
	Args: func(cms *cobra.Command, args []string) error {
		if (Add || Remove) && len(args) != 1 {
			return errors.New("expects exactly one argument")
		}
		if Add && Remove {
			return errors.New("add and remove cannot be combined")
		}
		if !(Add || Remove) && len(args) > 0 {
			return errors.New("too many arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.New()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Unable to read configuration.")
			os.Exit(1)
		}

		pg, err := postgres.New(*cfg.Postgres)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Unable to connect to Postgres:", err)
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
				_, _ = fmt.Fprintln(os.Stderr, "Unable to add admin:", err)
				os.Exit(1)
			}
			fmt.Printf("Added %s.\n", username)
		} else if Remove {
			username := args[0]
			err := us.DeleteUser(username)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Unable to remove admin:", err)
				os.Exit(1)
			}
			fmt.Printf("Removed %s.\n", username)
		} else {
			users, err := us.Users()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "unable to get users:", err)
				os.Exit(1)
			}

			if len(users) == 0 {
				_, _ = fmt.Println("No Shuttle Tracker administrators.")
				return
			}

			for _, user := range users {
				_, _ = fmt.Println(user.Username)
			}
		}
	},
}
