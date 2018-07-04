package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/bialas1993/go-crawler/crawler"
	"golang.org/x/net/html"
	"context"
	"github.com/bialas1993/go-crawler/filters"
	"github.com/joho/godotenv"
	"os"
	)

var pageUrl, logTimer, logLevel, auth = parseParams()

func init() {
	log.SetLevel(log.Level(int32(logLevel)))
	godotenv.Load()
}

func main() {
	var logger crawler.LogService
	nodeChan := make(chan *html.Node)
	ctx, cancel := context.WithCancel(context.Background())
	fm := filters.CreateManager()

	if len(os.Getenv(ELASTICSEARCH_LOGGER_HOST_ENV)) == 0 {
		logger = CreateLogger()
	} else {
		logger = CreateElasticLoggerService(crawler.ExtraDataMessage{Domain:pageUrl})
	}

	log.Printf("pageUrl=%s, logTimer=%d, logLevel=%d\n", pageUrl,  logTimer, logLevel)

	go func() {
		crawler.Run(pageUrl, &nodeChan, ctx, cancel, logger, auth)
	} ()

	for {
		select {
		case node := <- nodeChan:
			if node != nil {
				fm.Parse(node)
				continue
			}
			break
		case <-ctx.Done():
			close(nodeChan)
			return
		}
	}
}