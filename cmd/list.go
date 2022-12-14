package cmd

import (
	"github.com/carloscastrojumo/remindme/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var tags []string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List notes to the database",
	Long:  `Add new note to the database`,
	Run: func(cmd *cobra.Command, args []string) {
		notes, err := noteService.GetByTags(tags)
		if err != nil {
			color.Red("Error while getting notes by tags: %s", err)
		}
		output.Print(notes)
	},
}

func init() {
	listCmd.Flags().StringArrayVar(&tags, "tags", []string{}, "Tags to add to the note")
	listCmd.MarkFlagRequired("tags")
	rootCmd.AddCommand(listCmd)
}
