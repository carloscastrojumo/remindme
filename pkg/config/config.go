package config

import (
	"errors"
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
func InitConfig() error {
	viper.AddConfigPath(appDir)
	viper.SetConfigName("config")

	// create new folder if it doesn't exist
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("create config folder %s: %w", appDir, err)
	}

	err := viper.ReadInConfig()
	var notFound viper.ConfigFileNotFoundError
	if errors.As(err, &notFound) {
		fmt.Fprintln(os.Stderr, "Config file not found, creating one")
		return promptConfigFile()
	}
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}
	return nil
}

func promptConfigFile() error {
	storageType, err := prompt.ForString("What storage type do you want to use? (mongo, yaml) [yaml]")
	if err != nil {
		return err
	}
	if len(storageType) == 0 {
		storageType = "yaml"
	}

	viper.Set("storageType", storageType)

	switch storageType {
	case "mongo":
		for _, key := range []string{"host", "port", "database", "collection"} {
			value, err := prompt.ForString("Mongo " + key)
			if err != nil {
				return err
			}
			viper.Set("mongo."+key, value)
		}
	case "yaml":
		dataFilename, err := prompt.ForString("YAML file name (current directory: " + appDir + ") [data.yaml]")
		if err != nil {
			return err
		}
		if len(dataFilename) == 0 {
			dataFilename = "data.yaml"
		}
		viper.Set("yaml.name", appDir+"/"+dataFilename)
	}

	return saveConfigFile()
}

func saveConfigFile() error {
	configDir := xdg.Home + "/.config/remindme"
	viper.AddConfigPath(configDir)
	if err := viper.WriteConfigAs(configDir + "/config.yaml"); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}
	return nil
}

// GetNoteService returns a new note service
func GetNoteService() (*storage.NoteService, error) {
	config.StorageType = viper.GetString("storageType")

	switch config.StorageType {
	case "mongo":
		color.New(color.FgBlue).Fprintln(os.Stderr, "Using Mongo storage")
		var mongoConfig mongo.Config
		if err := viper.UnmarshalKey("mongo", &mongoConfig); err != nil {
			return nil, fmt.Errorf("read %s configuration: %w", config.StorageType, err)
		}

		config.StorageConfig = &mongoConfig

	case "yaml":
		color.New(color.FgBlue).Fprintln(os.Stderr, "Using YAML storage")
		var yamlConfig yaml.Config
		if err := viper.UnmarshalKey("yaml", &yamlConfig); err != nil {
			return nil, fmt.Errorf("read %s configuration: %w", config.StorageType, err)
		}

		config.StorageConfig = &yamlConfig

	default:
		return nil, fmt.Errorf("unsupported storage type %q, expected mongo or yaml", config.StorageType)
	}

	return initNoteService(config)
}

func initNoteService(storageConfig *storage.Config) (*storage.NoteService, error) {
	storeService, err := storage.GetStorage(storageConfig)
	if err != nil {
		return nil, fmt.Errorf("initialize %s storage: %w", storageConfig.StorageType, err)
	}
	return storage.NewNoteService(storeService), nil
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
