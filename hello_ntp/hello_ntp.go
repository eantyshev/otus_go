package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

func main() {
	t := time.Now()
	fmt.Println("Local time: ", t.Format(time.RFC3339))
	ntpServer := os.Args[1]
	tNtp, err := ntp.Time(ntpServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query ntp server: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("NTP time: ", tNtp.Format(time.RFC3339))
}
