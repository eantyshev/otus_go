package config

import (
	"fmt"
	logger2 "github.com/eantyshev/otus_go/calendar/internal/logger"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

func init() {
	flag.StringVar(&cfgFile, "config",
		"config.yaml",
		"config file (default is ./config.yaml)")
}

func SetupViper() {
	viper.SetConfigFile(cfgFile)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetEnvPrefix("calendar")
	viper.AutomaticEnv()
}

func Configure() {
	SetupViper()
	logger2.ConfigureLogger()
}
