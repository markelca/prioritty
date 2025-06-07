package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const CONF_DATABASE_PATH string = "database_path"
const CONF_LOG_FILE_PATH string = "log_file_path"
const CONF_DEFAULT_COMMAND string = "default_command"
const CONF_EDITOR string = "editor"

type Config struct {
	DatabasePath   string `mapstructure:"database_path"`
	LogFilePath    string `mapstructure:"log_file_path"`
	DefaultCommand string `mapstructure:"default_command"`
	Editor         string `mapstructure:"editor"`
}

var config *Config

func InitConfig(cfgFile string) error {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(configPath, "prioritty")

	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name "prioritty" (without extension)
		// viper.AddConfigPath(".")
		// viper.AddConfigPath(home)
		viper.AddConfigPath(configDir)

		viper.SetConfigName("prioritty")
		viper.SetConfigType("yaml")
	}

	setDefaults(configDir)

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("PRIORITTY")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createConfigFile(configDir); err != nil {
				return fmt.Errorf("Failed to create the config file: %w", err)
			}
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("unable to decode config: %w", err)
	}

	return nil
}

func createConfigFile(configDir string) error {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	// Create the config file with default values
	if err := viper.SafeWriteConfig(); err != nil {
		log.Fatalf("Error creating config file: %v", err)
	}
	return nil
}

func setDefaults(configDir string) {
	viper.SetDefault(CONF_DATABASE_PATH, filepath.Join(configDir, "prioritty.db"))
	viper.SetDefault(CONF_LOG_FILE_PATH, filepath.Join(configDir, "prioritty.log"))
	viper.SetDefault(CONF_DEFAULT_COMMAND, "tui")
	viper.SetDefault(CONF_EDITOR, "nano")
}
