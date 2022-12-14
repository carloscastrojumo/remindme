package cmd

import (
	"fmt"
	"os"

	"github.com/carloscastrojumo/remindme/pkg/storage"
	"github.com/carloscastrojumo/remindme/pkg/storage/mongo"
	"github.com/carloscastrojumo/remindme/pkg/storage/yaml"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func Execute() {
	setStorageConfig()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func setStorageConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Fprintf(os.Stderr, "Whoops. Config file '%s' not found", err)
			os.Exit(1)
		} else {
			// Config file was found but another error was produced
			fmt.Fprintf(os.Stderr, "Whoops. There was an error while reading your config file '%s'", err)
			os.Exit(1)
		}
	}

	storageType := viper.GetString("storageType")

	switch storageType {
	case "mongo":
		color.Blue("Using Mongo storage")
		var mongoConfig mongo.Config
		err := viper.UnmarshalKey("mongo", &mongoConfig)

		if err != nil {
			color.Red("Could not read %s configuration: '%s'", storageType, err)
			os.Exit(1)
		}

		storageConfig := &storage.StorageConfig{
			StorageType:   "mongo",
			StorageConfig: &mongoConfig,
		}

		initNoteService(storageConfig)
	case "yaml":
		color.Blue("Using YAML storage")
		var yamlConfig yaml.Config
		err := viper.UnmarshalKey("yaml", &yamlConfig)

		if err != nil {
			color.Red("Whoops. Could not read %s configuration: '%s'", storageType, err)
			os.Exit(1)
		}

		storageConfig := &storage.StorageConfig{
			StorageType:   "yaml",
			StorageConfig: &yamlConfig,
		}

		initNoteService(storageConfig)
	default:
		color.Red("No storage type found")
		os.Exit(1)
	}
}

func initNoteService(storageConfig *storage.StorageConfig) {
	storeService := storage.GetStorage(storageConfig)
	noteService = storage.NewNoteService(storeService)
}
