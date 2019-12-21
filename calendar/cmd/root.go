package cmd

import (
	"fmt"
	"github.com/eantyshev/otus_go/calendar/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Very basic calendar web service",
	Long: `
Calendar provides HTTP API for CRUD-equivalent
operations on appointments.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(logger.ConfigureLogger)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config",
		"config.yaml", "config file (default is ./config.yaml)")

	rootCmd.AddCommand(rpcServerCmd)
	rootCmd.AddCommand(rpcClientCmd)

	rpcClientCmd.Flags().StringVarP(&Args.Uuid, "uuid","u", "", "uuid of entity")
	rpcClientCmd.Flags().StringVar(&Args.RequestJson, "request-json", "", "entity info as json")
	rpcClientCmd.Flags().StringVar(&Args.Owner, "owner", "", "owner identity")
	rpcClientCmd.Flags().StringVar(&Args.Period, "period", "", "Period to list.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(cfgFile)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetEnvPrefix("calendar")
	viper.AutomaticEnv()
}
