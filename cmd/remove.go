package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove note from the database",
	Long:  `Remove note from the database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		if id != "" {
			if err := noteService.Remove(id); err != nil {
				return fmt.Errorf("delete note %s: %w", id, err)
			}
			color.Green("Note %s deleted", id)
		}

		if len(tags) > 0 {
			if err := noteService.RemoveByTags(tags); err != nil {
				return fmt.Errorf("delete notes by tags: %w", err)
			}
			color.Green("Notes with tags %s deleted", tags)
		}
		return nil
	},
}

func init() {
	removeCmd.Flags().String("id", "", "ID of the note to remove")
	removeCmd.Flags().StringSlice("tags", []string{}, "Remove all notes with any of these tags, comma separated or repeated")
	removeCmd.MarkFlagsOneRequired("id", "tags")
	rootCmd.AddCommand(removeCmd)
}
