package main

import (
	"flag"
	"log"

	"github.com/koron/nvgd/core"
)

var (
	config = flag.String("c", "", "configuration file")
)

func main() {
	flag.Parse()
	c, err := core.LoadConfig(*config)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	if err := core.Run(c); err != nil {
		log.Fatalf("failed to run server: %s", err)
	}
}
