package cmd

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/carloscastrojumo/remindme/pkg/storage"
	"github.com/carloscastrojumo/remindme/pkg/storage/mongo"
	"github.com/carloscastrojumo/remindme/pkg/storage/yaml"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
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
	configDir := xdg.Home + "/.config/remindme"

	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")

	// create new folder if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if os.Mkdir(configDir, 0755) != nil {
			color.Red("Error while creating config folder: %s", err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create one
			fmt.Fprintf(os.Stderr, "Config file not found, creating one")
			initConfigFile()
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

func initConfigFile() {
	storageType := PromptForString("What storage type do you want to use? (mongo, yaml)")

	viper.Set("storageType", storageType)

	switch storageType {
	case "mongo":
		viper.Set("mongo.host", PromptForString("Mongo host"))
		viper.Set("mongo.port", PromptForString("Mongo port"))
		viper.Set("mongo.database", PromptForString("Mongo database"))
		viper.Set("mongo.collection", PromptForString("Mongo collection"))
	case "yaml":
		viper.Set("yaml.name", PromptForString("YAML file name"))
	}

	saveConfigFile()
}

func PromptForString(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}

func saveConfigFile() {
	configDir := xdg.Home + "/.config/remindme"
	viper.AddConfigPath(configDir)
	viper.WriteConfigAs(configDir + "/config.yaml")
}

func initNoteService(storageConfig *storage.StorageConfig) {
	storeService := storage.GetStorage(storageConfig)
	noteService = storage.NewNoteService(storeService)
}
