package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/adrg/xdg"
	prompt "github.com/carloscastrojumo/remindme/pkg/prompt"
	"github.com/carloscastrojumo/remindme/pkg/storage"
	"github.com/carloscastrojumo/remindme/pkg/storage/mongo"
	"github.com/carloscastrojumo/remindme/pkg/storage/yaml"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var appDir = configDir()

// configDir honors an absolute $XDG_CONFIG_HOME, but keeps using
// ~/.config/remindme when a config file already exists there.
func configDir() string {
	legacy := filepath.Join(xdg.Home, ".config", "remindme")
	base := os.Getenv("XDG_CONFIG_HOME")
	if !filepath.IsAbs(base) {
		return legacy
	}
	if _, err := os.Stat(filepath.Join(legacy, "config.yaml")); err == nil {
		return legacy
	}
	return filepath.Join(base, "remindme")
}

var storageConfig any

// InitConfig initializes the configuration
func InitConfig() error {
	viper.AddConfigPath(appDir)
	viper.SetConfigName("config")

	// create new folder if it doesn't exist
	if err := os.MkdirAll(appDir, 0700); err != nil {
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
		viper.Set("yaml.name", filepath.Join(appDir, dataFilename))
	}

	return saveConfigFile()
}

func saveConfigFile() error {
	viper.SetConfigPermissions(0600)
	if err := viper.WriteConfigAs(filepath.Join(appDir, "config.yaml")); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}
	return nil
}

// GetNoteService returns a new note service
func GetNoteService() (*storage.NoteService, error) {
	storageType := viper.GetString("storageType")

	var store storage.NoteStorage
	switch storageType {
	case "mongo":
		color.New(color.FgBlue).Fprintln(os.Stderr, "Using Mongo storage")
		var mongoConfig mongo.Config
		if err := viper.UnmarshalKey("mongo", &mongoConfig); err != nil {
			return nil, fmt.Errorf("read %s configuration: %w", storageType, err)
		}
		storageConfig = &mongoConfig

		mongoStore, err := mongo.Initialize(&mongoConfig)
		if err != nil {
			return nil, fmt.Errorf("initialize %s storage: %w", storageType, err)
		}
		store = mongoStore

	case "yaml":
		color.New(color.FgBlue).Fprintln(os.Stderr, "Using YAML storage")
		var yamlConfig yaml.Config
		if err := viper.UnmarshalKey("yaml", &yamlConfig); err != nil {
			return nil, fmt.Errorf("read %s configuration: %w", storageType, err)
		}
		storageConfig = &yamlConfig

		yamlStore, err := yaml.Initialize(&yamlConfig)
		if err != nil {
			return nil, fmt.Errorf("initialize %s storage: %w", storageType, err)
		}
		store = yamlStore

	default:
		return nil, fmt.Errorf("unsupported storage type %q, expected mongo or yaml", storageType)
	}

	return storage.NewNoteService(store), nil
}

// GetConfig prints the current configuration to screen
func GetConfig() {
	color.Blue("Configuration file: %s\n", color.GreenString(viper.ConfigFileUsed()))
	switch c := storageConfig.(type) {
	case *yaml.Config:
		color.Blue("Data file: %s\n", color.GreenString(c.Name))
	case *mongo.Config:
		color.Blue("Host: %s\n", color.GreenString(c.Host))
		color.Blue("Port: %s\n", color.GreenString(strconv.Itoa(c.Port)))
		color.Blue("Database: %s\n", color.GreenString(c.Database))
		color.Blue("Collection: %s\n", color.GreenString(c.Collection))
	}
}
