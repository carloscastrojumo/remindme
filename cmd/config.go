package cmd

import (
	"github.com/carloscastrojumo/remindme/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Prints the current configuration file to screen",
	Long:  `Prints the current configuration file to screen`,
	Run: func(cmd *cobra.Command, args []string) {
		config.GetConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
