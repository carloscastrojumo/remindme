package cmd

import (
	"github.com/carloscastrojumo/remindme/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listTags = &cobra.Command{
	Use:   "tags",
	Short: "List all tags available",
	Long:  "List all tags available",
	Run: func(cmd *cobra.Command, args []string) {
		tags, err := noteService.GetTags()
		if err != nil {
			color.Red("Error while getting notes by tags: %s", err)
		}
		output.PrintTags(tags)
	},
}

func init() {
	listCmd.AddCommand(listTags)
}
