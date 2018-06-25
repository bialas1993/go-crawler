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

	logChan := make(chan crawler.LogMessage)
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
			if logMsg.Level() != crawler.LOG_MESSAGE_CLOSE {
				switch logMsg.Level() {
				case crawler.LOG_MESSAGE_INFO:
					log.Info(logMsg.GetMessage())
				case crawler.LOG_MESSAGE_WARN:
					log.Warn(logMsg.GetMessage())
				case crawler.LOG_MESSAGE_ERROR:
					log.Error(logMsg.GetMessage())
				case crawler.LOG_MESSAGE_DEBUG:
					log.Debug(logMsg.GetMessage())
				}
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