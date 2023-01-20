package cmd

import (
	"github.com/carloscastrojumo/remindme/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configuration",
	Long:  `Helps manage application configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// var searchLocations []string

		// cmd.Flags().Visit(func(f *pflag.Flag) {
		// 	searchLocations = append(searchLocations, f.Name)
		// })

		// if len(searchLocations) == 0 {
		// 	searchLocations = append(searchLocations, "command", "description", "tags")
		// }

		// if len(args) < 1 {
		// 	return errors.New("the word-to-search argument is required")
		// }

		// notes, err := noteService.Search(args[0], searchLocations)
		// if err != nil {
		// 	color.Red("Error while getting notes by tags: %s", err)
		// }
		// output.Print(notes)
		config.GetConfig()
	},
}

func init() {
	// searchCmd.Flags().BoolP("tags", "t", false, "Search in tags")
	// searchCmd.Flags().BoolP("command", "c", false, "Search in commands")
	// searchCmd.Flags().BoolP("description", "d", false, "Search in description")
	rootCmd.AddCommand(configCmd)
}
