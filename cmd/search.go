package cmd

import (
	"errors"
	"strings"

	"github.com/carloscastrojumo/remindme/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var searchCmd = &cobra.Command{
	Use:   "search [word-to-search] [flags]",
	Short: "Search the notes in the database",
	Long:  `Search the notes in the database, by default it searches in tags, commands, and descriptions. If a flag is specified it will only search in those.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var searchLocations []string

		cmd.Flags().Visit(func(f *pflag.Flag) {
			searchLocations = append(searchLocations, f.Name)
		})

		if len(searchLocations) == 0 {
			searchLocations = append(searchLocations, "command", "description", "tags")
		}

		if len(args) < 1 {
			return errors.New("the word-to-search argument is required")
		}

		words := strings.Split(args[0], " ")
		notes, err := noteService.Search(words, searchLocations)
		if err != nil {
			color.Red("Error while getting notes by tags: %s", err)
		}
		output.Print(notes)

		return nil
	},
}

func init() {
	searchCmd.Flags().BoolP("tags", "t", false, "Search in tags")
	searchCmd.Flags().BoolP("command", "c", false, "Search in commands")
	searchCmd.Flags().BoolP("description", "d", false, "Search in description")
	rootCmd.AddCommand(searchCmd)
}
