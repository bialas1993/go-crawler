package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/bialas1993/go-crawler/crawler"
	"golang.org/x/net/html"
	"context"
	"github.com/bialas1993/go-crawler/filters"
	)

var pageUrl, depth, logTimer, logLevel = parseParams()

func init() {
	log.SetLevel(log.Level(int32(logLevel)))
}

func main() {
	nodeChan := make(chan *html.Node)
	ctx, cancel := context.WithCancel(context.Background())
	fm := filters.NewManager()
	logger := CreateLogger()
	log.Printf("pageUrl=%s, depth=%d, logTimer=%d, logLevel=%d", pageUrl, depth, logTimer, logLevel)

	go func() {
		crawler.Run(pageUrl, &nodeChan, ctx, cancel, logger)
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