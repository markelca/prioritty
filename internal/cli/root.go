/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"log"

	"github.com/markelca/prioritty/internal/config"
	"github.com/markelca/prioritty/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pt",
	Short: "Task list and productivity tool in the terminal",
	Long: `A Terminal User Interface (TUI) and CLI application for managing your tasks. Focused on:
	- Good looks
	- Performance
	- Nice defaults
	- Customization
	- Autocompletion support

🚧 Disclaimer: This project is still under development.`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultCommand := viper.GetString("default_command")
		if subCmd, _, err := cmd.Find([]string{defaultCommand}); err == nil {
			subCmd.Run(cmd, args)
		} else {
			log.Printf("Error - The default configured command doesn't exist (%s)", defaultCommand)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

var cfgFile string

func init() {
	cobra.OnInitialize(func() {
		config.InitConfig(cfgFile)
		if err := logger.InitLogger(); err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}
	})
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.Flags().Bool("demo", false, "Use a temporal demo database with predefined values")
	viper.BindPFlag("demo", rootCmd.Flags().Lookup("demo"))
}
