package cmd

import (
	"fmt"

	"github.com/carloscastrojumo/remindme/internal/output"
	"github.com/spf13/cobra"
)

var listAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all notes in the database",
	Long:  "List all notes in the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		notes, err := noteService.GetAll()
		if err != nil {
			return fmt.Errorf("get all notes: %w", err)
		}
		output.Print(notes)
		return nil
	},
}

func init() {
	listCmd.AddCommand(listAllCmd)
}
