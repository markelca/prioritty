/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"os"

	"github.com/markelca/prioritty/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pt",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultCommand := viper.GetString("default_command")
		if subCmd, _, err := cmd.Find([]string{defaultCommand}); err == nil {
			subCmd.Run(cmd, args)
		} else {
			fmt.Printf("Error - The default configured command doesn't exist (%s)", defaultCommand)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	cobra.OnInitialize(func() {
		config.InitConfig(cfgFile)
	})
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")

	rootCmd.PersistentFlags().Bool("demo", false, "Populate for demo values for showcase")
	viper.BindPFlag("demo", rootCmd.PersistentFlags().Lookup("demo"))
	viper.SetDefault("demo", "false")
}
