package cmd

import (
	"github.com/spf13/viper"

	"github.com/eantyshev/otus_go/calendar/hello"
	"github.com/spf13/cobra"
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Start primitive hello-world http server",
	Run: func(cmd *cobra.Command, args []string) {
		hello.Server(viper.GetString("http_listen"))
	},
}
