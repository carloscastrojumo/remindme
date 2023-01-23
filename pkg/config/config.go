package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/adrg/xdg"
	prompt "github.com/carloscastrojumo/remindme/pkg/prompt"
	"github.com/carloscastrojumo/remindme/pkg/storage"
	"github.com/carloscastrojumo/remindme/pkg/storage/mongo"
	"github.com/carloscastrojumo/remindme/pkg/storage/yaml"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var appDir = xdg.Home + "/.config/remindme"

var config = &storage.Config{}

// InitConfig initializes the configuration
func InitConfig() {
	viper.AddConfigPath(appDir)
	viper.SetConfigName("config")

	// create new folder if it doesn't exist
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		if os.Mkdir(appDir, 0755) != nil {
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
		dataFilename := prompt.ForString("YAML file name (current directory: " + appDir + ")")
		viper.Set("yaml.name", appDir+"/"+dataFilename)
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
	config.StorageType = viper.GetString("storageType")

	switch config.StorageType {
	case "mongo":
		color.Blue("Using Mongo storage")
		var mongoConfig mongo.Config
		err := viper.UnmarshalKey("mongo", &mongoConfig)

		if err != nil {
			color.Red("Could not read %s configuration: '%s'", config.StorageType, err)
			os.Exit(1)
		}

		config.StorageConfig = &mongoConfig

	case "yaml":
		color.Blue("Using YAML storage")
		var yamlConfig yaml.Config
		err := viper.UnmarshalKey("yaml", &yamlConfig)

		if err != nil {
			color.Red("Whoops. Could not read %s configuration: '%s'", config.StorageType, err)
			os.Exit(1)
		}

		config.StorageConfig = &yamlConfig

	default:
		color.Red("No storage type found")
		os.Exit(1)
	}

	return initNoteService(config)
}

func initNoteService(storageConfig *storage.Config) *storage.NoteService {
	storeService := storage.GetStorage(storageConfig)
	return storage.NewNoteService(storeService)
}

// GetConfig prints the current configuration to screen
func GetConfig() {
	color.Blue("Configuration file: %s\n", color.GreenString(viper.ConfigFileUsed()))
	switch config.StorageType {
	case "yaml":
		color.Blue("Data file: %s\n", color.GreenString(config.StorageConfig.(*yaml.Config).Name))
	case "mongo":
		color.Blue("Host: %s\n", color.GreenString(config.StorageConfig.(*mongo.Config).Host))
		color.Blue("Port: %s\n", color.GreenString(strconv.Itoa(config.StorageConfig.(*mongo.Config).Port)))
		color.Blue("Database: %s\n", color.GreenString(config.StorageConfig.(*mongo.Config).Database))
		color.Blue("Collection: %s\n", color.GreenString(config.StorageConfig.(*mongo.Config).Collection))
	}
}
