package main

import (
	"fmt"
	"os"
	"flag"
)

func usage() {
	fmt.Printf("usage: go-crawler http:/your-domain.com/ --depth=5 \n")
	flag.PrintDefaults()
	os.Exit(2)
}

func parseParams() (string, int) {
	flag.Usage = usage
	depth := flag.Int("depth", 0, "Depth to finding pages")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		usage()
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	*depth += 1

	return string(args[0]), *depth
}
