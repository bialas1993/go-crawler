package main

import (
	"bufio"
	"fmt"
	"os"
	"flag"
	"strings"
	"github.com/bialas1993/go-crawler/crawler"
	"github.com/howeyc/gopass"
)

func usage() {
	fmt.Printf("usage: go-crawler -pageUrl=http://your-domain.com/ -timer=1 -auth=1\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func parseParams() (string, int, int, crawler.AuthCredentials) {
	var logTimer, logLevel, authEnabled int
	var page string
	auth := crawler.AuthCredentials {
		Enabled: false,
	}

	flag.Usage = usage
	flag.StringVar(&page, "page", "", "Url to parse pageUrl")
	flag.IntVar(&logTimer, "timer", 1, "Timer to log performance")
	flag.IntVar(&logLevel, "level", crawler.LOG_MESSAGE_INFO, "Logging level")
	flag.IntVar(&authEnabled, "auth", crawler.AUTH_DISABLE, "Auth enabled")

	flag.Parse()

	if len(page) < 1 {
		usage()
		fmt.Println("Please specify start pageUrl")
		os.Exit(1)
	}

	if authEnabled == crawler.AUTH_ENABLE {
		auth = credentials()
	}

	return page, logTimer, logLevel, auth
}

func credentials() crawler.AuthCredentials {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, _ := gopass.GetPasswd()
	password := string(bytePassword)

	return crawler.AuthCredentials{
		Enabled: true,
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}
}