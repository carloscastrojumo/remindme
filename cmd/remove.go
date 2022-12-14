package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var id string

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Add new note to the database",
	Long:  `Add new note to the database`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := noteService.Remove(id); err != nil {
			panic(err)
		}
		color.Green("Note deleted successfully")
	},
}

func init() {
	removeCmd.Flags().StringVar(&id, "id", "", "ID of the note to remove")
	removeCmd.Flags().StringArrayVar(&tags, "tags", []string{}, "Remove all notes from tags")
	rootCmd.AddCommand(removeCmd)
}
