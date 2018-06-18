package main

import (
	"fmt"
	"os"
	"flag"
)

func usage() {
	fmt.Printf("usage: go-crawler --pageUrl=http://your-domain.com/ --depth=5 --timer=1\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func parseParams() (string, int, int, int) {
	var depth, logTimer int
	var page string

	flag.Usage = usage
	flag.StringVar(&page, "page", "", "Url to parse pageUrl")
	flag.IntVar(&depth, "depth", 0, "Depth to finding pages")
	flag.IntVar(&logTimer, "timer", 1, "Timer to log performance")

	flag.Parse()

	if len(page) < 1 {
		usage()
		fmt.Println("Please specify start pageUrl")
		os.Exit(1)
	}

	depth += 1

	return page, depth, logTimer, LOG_LEVEL_DEBUG
}


func selectDebugLevel() int {
	for i := 0; i < len(os.Args);  {
		switch os.Args[i] {
		case "-v":
			return LOG_LEVEL_WARNING
		case "-vv":
			return LOG_LEVEL_ERROR
		case "-vvv":
			return LOG_LEVEL_DEBUG
		}
	}

	return LOG_LEVEL_DEBUG
}