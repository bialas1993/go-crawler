package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/bialas1993/go-crawler/crawler"
		"golang.org/x/net/html"
	"context"
	"github.com/bialas1993/go-crawler/filters"
)

func main() {
	pageUrl, depth, logTimer, logLevel := parseParams()
	logChan := make(chan string)
	nodeChan := make(chan *html.Node)
	ctx, cancel := context.WithCancel(context.Background())
	fm := filters.NewManager()

	log.Printf("pageUrl=%s, depth=%d, logTimer=%d, logLevel=%d", pageUrl, depth, logTimer, logLevel)

	go func() {
		crawler.Run(pageUrl, &nodeChan, &logChan, ctx, cancel)
	} ()

	for {
		select {
		case logMsg := <-logChan:
			if logMsg != crawler.CLOSE_LOGGER_MESSAGE {
				log.Println(logMsg)
				continue
			}
			break
		case node := <- nodeChan:
			if node != nil {
				fm.Parse(node)
				continue
			}
			break
		case <-ctx.Done():
			close(logChan)
			close(nodeChan)
			return
		}
	}
}