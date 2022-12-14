package cmd

import (
	"fmt"

	"github.com/carloscastrojumo/remindme/pkg/storage"
	"github.com/spf13/cobra"
)

var note storage.Note

func init() {
	addCmd.Flags().StringArrayVar(&note.Tags, "tags", []string{}, "Tags to add to the note")
	addCmd.Flags().StringVar(&note.Command, "command", "", "Command to add to the note")
	addCmd.Flags().StringVar(&note.Description, "description", "", "Description to add to the note")
	addCmd.MarkFlagRequired("tags")
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new note to the database",
	Long:  `Add new note to the database`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := noteService.Add(note); err != nil {
			panic(err)
		}
		fmt.Println("Note added successfully")
	},
}
