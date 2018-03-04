package main

import (
	"fmt"
	"os"
	"flag"
)

func usage() {
	fmt.Printf("usage: go-crawler http:/your-domain.com/\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func validateParams() []string {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		usage()
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	return args
}
