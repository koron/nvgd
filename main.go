package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/core"
	"github.com/koron/nvgd/internal/version"
)

var (
	configOpt = flag.String("c", "nvgd.conf.yml", "configuration file")
	pprofAddr = flag.String("pprofaddr", "", "address for pprof server")
	verFlag   = flag.Bool("version", false, "show version")
)

func main() {
	flag.Parse()
	if *verFlag {
		fmt.Println(version.Version)
		return
	}
	c, err := config.LoadConfig(*configOpt)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	if *pprofAddr != "" {
		log.Printf("start pprof on %s", *pprofAddr)
		go func() {
			log.Println(http.ListenAndServe(*pprofAddr, nil))
		}()
	}
	if err := core.Run(c); err != nil {
		log.Fatalf("failed to run server: %s", err)
	}
}
