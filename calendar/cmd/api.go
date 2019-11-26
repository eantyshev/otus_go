package cmd

import (
	"github.com/eantyshev/otus_go/calendar/api"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// apiCmd starts REST API service
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start REST API http server",
	Run: func(cmd *cobra.Command, args []string) {
		api.Server(viper.GetString("http_listen"))
	},
}
