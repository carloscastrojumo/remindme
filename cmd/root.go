package cmd

import (
	"fmt"
	"os"

	"github.com/carloscastrojumo/remindme/internal/config"
	"github.com/carloscastrojumo/remindme/internal/storage"
	"github.com/spf13/cobra"
)

var noteService *storage.NoteService

var rootCmd = &cobra.Command{
	Use:   "rmm",
	Short: "remindme - a simple CLI to remind you about notes",
	Long:  `remindme - a simple CLI to remind you about notes`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		if !needsStorage(cmd) {
			return nil
		}
		if err := config.InitConfig(); err != nil {
			return err
		}
		var err error
		noteService, err = config.GetNoteService()
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remindme - a simple CLI to remind you about notes")
	},
}

func needsStorage(cmd *cobra.Command) bool {
	for c := cmd; c != nil; c = c.Parent() {
		switch c.Name() {
		case "help", "completion", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
			return false
		}
	}
	return true
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
