package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const KEY_DATABASE_PATH string = "database_path"
const KEY_LOG_FILE_PATH string = "log_file_path"
const KEY_DEFAULT_COMMAND string = "default_command"

type Config struct {
	DatabasePath   string `mapstructure:"database_path"`
	LogFilePath    string `mapstructure:"log_file_path"`
	DefaultCommand string `mapstructure:"default_command"`
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
	fmt.Println("Config file not found, creating with defaults...")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	// Create the config file with default values
	if err := viper.SafeWriteConfig(); err != nil {
		log.Fatalf("Error creating config file: %v", err)
	}
	fmt.Println("Config file created successfully")
	return nil
}

func setDefaults(configDir string) {
	viper.SetDefault(KEY_DATABASE_PATH, filepath.Join(configDir, "prioritty.db"))
	viper.SetDefault(KEY_LOG_FILE_PATH, filepath.Join(configDir, "prioritty.log"))
	viper.SetDefault(KEY_DEFAULT_COMMAND, "ls")
}
