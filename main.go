package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/bialas1993/go-crawler/crawler"
	)

func main() {
	pageUrl, depth, logTimer, logLevel := parseParams()
	logChan := make(chan string)

	log.Printf("pageUrl=%s, depth=%d, logTimer=%d, logLevel=%d", pageUrl, depth, logTimer, logLevel)

	go func() {
		crawler.Run(pageUrl, &logChan)
	} ()

	for {
		logMsg := <-logChan
		if logMsg != crawler.CLOSE_LOGGER_MESSAGE {
			log.Println(logMsg)
			continue
		}

		break
	}

	defer close(logChan)
}