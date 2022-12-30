package cmd

import (
	"github.com/carloscastrojumo/remindme/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List notes to the database",
	Long:    `Add new note to the database`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetStringArray("tags")
		id, _ := cmd.Flags().GetString("id")

		if id != "" {
			if note, err := noteService.Get(id); err != nil {
				color.Red("Error: %s", err)
			} else {
				output.Print(note)
			}
		}

		if len(tags) > 0 {
			notes, err := noteService.GetByTags(tags)
			if err != nil {
				color.Red("Error while getting notes by tags: %s", err)
			}
			output.Print(notes)
		}

		if len(tags) == 0 && id == "" {
			notes, err := noteService.GetAll()
			if err != nil {
				color.Red("Error while getting all notes: %s", err)
			}
			output.Print(notes)
		}
	},
}

func init() {
	listCmd.Flags().StringArray("tags", []string{}, "Tags to add to the note")
	listCmd.Flags().String("id", "", "ID of the note")
	rootCmd.AddCommand(listCmd)
}
