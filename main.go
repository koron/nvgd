package main

import (
	"flag"
	"log"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/core"
)

var (
	configOpt = flag.String("c", "nvgd.conf.yml", "configuration file")
)

func main() {
	flag.Parse()
	c, err := config.LoadConfig(*configOpt)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	if err := core.Run(c); err != nil {
		log.Fatalf("failed to run server: %s", err)
	}
}
