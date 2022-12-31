package cmd

import (
	"fmt"

	"github.com/carloscastrojumo/remindme/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [word-to-search] [flags]",
	Short: "Search the notes in the database",
	Long:  `Search the notes in the database, by default it searches in tags, commands, and descriptions. If a flag is specified it will only search in those.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		fmt.Println(cmd.Flags().Changed("tags"))
		notes, err := noteService.Search(args[0], nil)
		if err != nil {
			color.Red("Error while getting notes by tags: %s", err)
		}
		output.Print(notes)
	},
}

func init() {
	searchCmd.Flags().BoolP("tags", "t", false, "Search in tags")
	searchCmd.Flags().BoolP("command", "c", false, "Search in commands")
	searchCmd.Flags().BoolP("description", "d", false, "Search in description")
	rootCmd.AddCommand(searchCmd)
}
