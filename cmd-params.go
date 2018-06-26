package main

import (
	"fmt"
	"os"
	"flag"
	"github.com/bialas1993/go-crawler/crawler"
)

const (
	LOG_LEVEL_INFO = iota
	LOG_LEVEL_WARNING
	LOG_LEVEL_ERROR
	LOG_LEVEL_DEBUG
)

func usage() {
	fmt.Printf("usage: go-crawler --pageUrl=http://your-domain.com/ --depth=5 --timer=1\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func parseParams() (string, int, int, int) {
	var depth, logTimer, logLevel int
	var page string

	flag.Usage = usage
	flag.StringVar(&page, "page", "", "Url to parse pageUrl")
	flag.IntVar(&depth, "depth", 0, "Depth to finding pages")
	flag.IntVar(&logTimer, "timer", 1, "Timer to log performance")
	flag.IntVar(&logLevel, "level", crawler.LOG_MESSAGE_INFO, "Logging level")

	flag.Parse()

	if len(page) < 1 {
		usage()
		fmt.Println("Please specify start pageUrl")
		os.Exit(1)
	}

	depth += 1

	return page, depth, logTimer, logLevel
}