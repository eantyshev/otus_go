package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"time"
)

var (
	host, port    string
	timeoutString string
)

func init() {
	pflag.StringVar(&timeoutString, "timeout",
		"60s", "Timeout on connect/send/receive.")
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s <host> <port> [--timeout n]\n", os.Args[0])
		pflag.PrintDefaults()
	}
}

func main() {
	pflag.Parse()
	if pflag.NArg() < 2 {
		pflag.Usage()
		os.Exit(1)
	}
	host = pflag.Arg(0)
	port = pflag.Arg(1)

	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tcpTalk(host, port, timeout); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
