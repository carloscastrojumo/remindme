package cmd

import (
	"github.com/carloscastrojumo/remindme/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// var tags []string

var listAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all notes in the database",
	Long:  "List all notes in the database",
	Run: func(cmd *cobra.Command, args []string) {
		notes, err := noteService.GetAll()
		if err != nil {
			color.Red("Error while getting notes by tags: %s", err)
		}
		output.Print(notes)
	},
}

func init() {
	listCmd.AddCommand(listAllCmd)
}
