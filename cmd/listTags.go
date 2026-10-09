package cmd

import (
	"fmt"

	"github.com/carloscastrojumo/remindme/internal/output"
	"github.com/spf13/cobra"
)

var listTags = &cobra.Command{
	Use:   "tags",
	Short: "List all tags available",
	Long:  "List all tags available",
	RunE: func(cmd *cobra.Command, args []string) error {
		tags, err := noteService.GetTags()
		if err != nil {
			return fmt.Errorf("get tags: %w", err)
		}
		output.PrintTags(tags)
		return nil
	},
}

func init() {
	listCmd.AddCommand(listTags)
}
