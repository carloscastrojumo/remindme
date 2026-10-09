package cmd

import (
	"fmt"

	prompt "github.com/carloscastrojumo/remindme/pkg/prompt"
	"github.com/carloscastrojumo/remindme/pkg/storage"
	"github.com/spf13/cobra"
)

var note storage.Note

func init() {
	addCmd.Flags().StringSliceVar(&note.Tags, "tags", []string{}, "Tags to add to the note, comma separated or repeated")
	addCmd.Flags().StringVar(&note.Command, "command", "", "Command to add to the note")
	addCmd.Flags().StringVar(&note.Description, "description", "", "Description to add to the note")
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new note to the database",
	Long:  `Add new note to the database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// if the user didn't provide any flags, we prompt for the note
		if note.Command == "" && note.Description == "" && len(note.Tags) == 0 {
			var err error
			if note, err = promptNote(); err != nil {
				return err
			}
		}

		if err := noteService.Add(note); err != nil {
			return fmt.Errorf("add note: %w", err)
		}

		fmt.Println("Note added successfully")
		return nil
	},
}

func promptNote() (storage.Note, error) {
	var note storage.Note
	var err error
	if note.Command, err = prompt.ForRequiredString("Command"); err != nil {
		return note, err
	}
	if note.Description, err = prompt.ForString("Description"); err != nil {
		return note, err
	}
	if note.Tags, err = prompt.ForRequiredStringArray("Tags"); err != nil {
		return note, err
	}
	return note, nil
}
