package cmd

import (
	"fmt"
	"os"

	"github.com/carloscastrojumo/remindme/pkg/config"
	"github.com/carloscastrojumo/remindme/pkg/storage"
	"github.com/spf13/cobra"
)

var noteService *storage.NoteService

var rootCmd = &cobra.Command{
	Use:   "remindme",
	Short: "remindme - a simple CLI to remind you about notes",
	Long: `remindme - a simple CLI to remind you about notes
   
One can use stringer to modify or inspect strings straight from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remindme - a simple CLI to remind you about notes")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	config.InitConfig()
	noteService = config.GetNoteService()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
