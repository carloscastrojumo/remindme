package config

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
	prompt "github.com/carloscastrojumo/remindme/pkg/prompt"
	"github.com/carloscastrojumo/remindme/pkg/storage"
	"github.com/carloscastrojumo/remindme/pkg/storage/mongo"
	"github.com/carloscastrojumo/remindme/pkg/storage/yaml"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// InitConfig initializes the configuration
func InitConfig() {
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
			promptConfigFile()
		} else {
			// Config file was found but another error was produced
			fmt.Fprintf(os.Stderr, "Whoops. There was an error while reading your config file '%s'", err)
			os.Exit(1)
		}
	}
}

func promptConfigFile() {
	storageType := prompt.ForString("What storage type do you want to use? (mongo, yaml)")

	viper.Set("storageType", storageType)

	switch storageType {
	case "mongo":
		viper.Set("mongo.host", prompt.ForString("Mongo host"))
		viper.Set("mongo.port", prompt.ForString("Mongo port"))
		viper.Set("mongo.database", prompt.ForString("Mongo database"))
		viper.Set("mongo.collection", prompt.ForString("Mongo collection"))
	case "yaml":
		viper.Set("yaml.name", prompt.ForString("YAML file name"))
	}

	saveConfigFile()
}

func saveConfigFile() {
	configDir := xdg.Home + "/.config/remindme"
	viper.AddConfigPath(configDir)
	viper.WriteConfigAs(configDir + "/config.yaml")
}

// GetNoteService returns a new note service
func GetNoteService() *storage.NoteService {
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

		return initNoteService(storageConfig)
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

		return initNoteService(storageConfig)
	default:
		color.Red("No storage type found")
		os.Exit(1)
	}

	return nil
}

func initNoteService(storageConfig *storage.StorageConfig) *storage.NoteService {
	storeService := storage.GetStorage(storageConfig)
	return storage.NewNoteService(storeService)
}
