package cli

import (
	"fmt"

	"github.com/markelca/prioritty/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Long:  `Display the current configuration values being used by prioritty.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		fmt.Println("Current Configuration:")
		fmt.Println("=====================")
		fmt.Printf("Config file: %s\n\n", viper.ConfigFileUsed())

		fmt.Printf("Database Path: %s\n", cfg.DatabasePath)
		fmt.Printf("Log File Path: %s\n", cfg.LogFilePath)
		fmt.Printf("Default Command: %s\n", cfg.DefaultCommand)
	},
}
