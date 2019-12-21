package cmd

import (
	"github.com/eantyshev/otus_go/calendar/rpc/server"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// rpcServerCmd starts GRPC API service
var rpcServerCmd = &cobra.Command{
	Use:   "rpc_server",
	Short: "Start GRPC API http server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Server(
			viper.GetString("http_listen"),
			viper.GetString("storage.pg.dsn"),
		)
	},
}
