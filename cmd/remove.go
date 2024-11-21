package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove note from the database",
	Long:  `Remove note from the database`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		tags, _ := cmd.Flags().GetStringArray("tags")

		if id != "" {
			if err := noteService.Remove(id); err != nil {
				color.Red("Error: %s", err)
			} else {
				color.Green("Note %s deleted", id)
			}
		}

		if len(tags) > 0 {
			if err := noteService.RemoveByTags(tags); err != nil {
				color.Red("Error while deleting notes by tags: %s", err)
			} else {
				color.Green("Notes with tags %s deleted", tags)
			}
		}
	},
}

func init() {
	removeCmd.Flags().String("id", "", "ID of the note to remove")
	removeCmd.Flags().StringArray("tags", []string{}, "Remove all notes from tags")
	rootCmd.AddCommand(removeCmd)
}
