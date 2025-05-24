package cli

import (
	"fmt"

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
		fmt.Printf("Database Path: %s\n", viper.Get("database_path"))
		fmt.Printf("Log File Path: %s\n", viper.Get("log_file_path"))
		fmt.Printf("Default Command: %s\n", viper.Get("default_command"))
	},
}
