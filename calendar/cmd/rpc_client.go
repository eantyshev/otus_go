package cmd

import (
	"github.com/eantyshev/otus_go/calendar/rpc/client"
	"github.com/spf13/viper"
	"log"

	"github.com/spf13/cobra"
)

var Args client.CallArgs

// rpcServerCmd starts GRPC API service
var rpcClientCmd = &cobra.Command{
	Use:   "rpc_client <action>",
	Short: "Execute GRPC call",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Args.Action = args[0]
		err := client.RpcCall(
			viper.GetString("http_listen"),
			viper.GetDuration("client.timeout"),
			Args,
			)
		if err != nil {
			log.Fatal(err)
		}
	},
}