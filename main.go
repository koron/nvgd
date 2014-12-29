package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/koron/nvd/night"
)

var config = flag.String("c", "nvd-conf.json", "configuration file")
var verbose = flag.Bool("v", false, "verbose message")
var help = flag.Bool("h", false, "show help message")

func showHelp() {
	fmt.Fprintln(os.Stderr, `USAGE: nvd [OPTIONS]

OPTIONS:`)
	flag.PrintDefaults()
}

func getLogger(v bool) *log.Logger {
	if v {
		return log.New(os.Stderr, "", log.LstdFlags)
	}
	return nil
}

func main() {
	flag.Parse()
	if *help {
		showHelp()
		return
	}
	if err := night.Run(*config, getLogger(*verbose)); err != nil {
		panic(err)
	}
}
