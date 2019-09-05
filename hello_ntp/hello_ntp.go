package main

import (
    "os"
    "fmt"
    "time"
    "github.com/beevik/ntp"
)

func main() {
    t := time.Now()
    fmt.Println("Local time: ", t.Format(time.RFC3339))
    ntp_server := os.Args[1]
    t_ntp, err := ntp.Time(ntp_server)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to query ntp server: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("NTP time: ", t_ntp.Format(time.RFC3339))
}
