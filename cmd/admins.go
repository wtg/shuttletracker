package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(adminsCmd)
}

var adminsCmd = &cobra.Command{
	Use:   "admins",
	Short: "Manage Shuttle Tracker administrators",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add")
	},
}
