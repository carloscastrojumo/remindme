package cmd

import (
	"fmt"

	"github.com/carloscastrojumo/remindme/pkg/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List notes from the database",
	Long:    `List notes from the database, optionally filtered by ID or tags`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tags, _ := cmd.Flags().GetStringSlice("tags")
		id, _ := cmd.Flags().GetString("id")

		if id != "" {
			note, err := noteService.Get(id)
			if err != nil {
				return err
			}
			output.Print([]interface{}{note})
		}

		if len(tags) > 0 {
			notes, err := noteService.GetByTags(tags)
			if err != nil {
				return fmt.Errorf("get notes by tags: %w", err)
			}
			output.Print(notes)
		}

		if len(tags) == 0 && id == "" {
			notes, err := noteService.GetAll()
			if err != nil {
				return fmt.Errorf("get all notes: %w", err)
			}
			output.Print(notes)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().StringSlice("tags", []string{}, "Filter notes by tag, comma separated or repeated")
	listCmd.Flags().String("id", "", "ID of the note")
	rootCmd.AddCommand(listCmd)
}
