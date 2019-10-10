package main

import (
	"flag"
	"log"
)

var pathFrom, pathTo string
var offset, limit int64

func init() {
	flag.StringVar(&pathFrom, "from", "", "file to read from")
	flag.StringVar(&pathTo, "to", "", "file to write to")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
	flag.Int64Var(&limit, "limit", 0, "limit the copied bytes")
}

func main() {
	flag.Parse()
	if err := copyData(pathFrom, pathTo, offset, limit); err != nil {
		log.Fatalf("Failed to copy: %s", err)
	}
}
