package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DatabasePath   string `mapstructure:"database_path"`
	LogFilePath    string `mapstructure:"log_file_path"`
	DefaultCommand string `mapstructure:"default_command"`
}

var GlobalConfig *Config

func InitConfig(cfgFile string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {

		configDir := filepath.Join(home, ".config", "prioritty")

		// Search config in home directory with name "prioritty" (without extension)
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath(configDir)
		viper.SetConfigName("prioritty")
		viper.SetConfigType("yaml")
	}

	// Set defaults
	setDefaults(home)

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("PRIORITTY")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("No config file found, using defaults")
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}

	// Unmarshal config into struct
	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("unable to decode config: %w", err)
	}

	return nil
}

func setDefaults(homeDir string) {
	// Database defaults
	viper.SetDefault("database_path", filepath.Join(homeDir, "share", "prioritty", "prioritty.db"))

	// Log file defaults
	viper.SetDefault("log_file_path", filepath.Join(homeDir, "config", "prioritty", "prioritty.logs"))

	// Default command to run when no subcommand is specified
	viper.SetDefault("default_command", "tui")
}

func GetConfig() *Config {
	if GlobalConfig == nil {
		panic("Config not initialized. Call InitConfig() first.")
	}
	return GlobalConfig
}
