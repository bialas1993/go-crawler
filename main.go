package main

import (
	"time"
	"sync"
	log "github.com/sirupsen/logrus"
	"net/url"
	)

const (
	LOG_LEVEL_INFO = iota
	LOG_LEVEL_WARNING
	LOG_LEVEL_ERROR
	LOG_LEVEL_DEBUG
	CONNECTIONS_LIMIT = 10
)

var tokens = make(chan struct{}, CONNECTIONS_LIMIT)
var httpErrors = make(chan *HttpError)
var mux sync.Mutex
var pageUrl string
var logTimer, depth, logLevel int
var page *url.URL

func crawl(url string) []string {
	tokens <- struct{}{}
	list, err := Extract(url)
	<-tokens

	if err != nil {
		err, ok := err.(*HttpError); if ok {
			go func (err *HttpError) {
				httpErrors <- err
			}(err)
		}
	}

	return list
}

func main() {
	var n int
	var pagesSeen, pageSeenLast uint16
	var timer *time.Ticker

	pageUrl, depth, logTimer, logLevel = parseParams()
	page, _ = url.Parse(pageUrl)
	worklist := make(chan []string)

	log.Printf("pageUrl=%s, depth=%d, logTimer=%d, logLevel=%d", pageUrl, depth, logTimer, logLevel)

	if logTimer > 0 {
		timer = time.NewTicker( time.Duration(logTimer) * time.Second )
	} else {
		timer = time.NewTicker( time.Second )
		timer.Stop()
	}

	n++
	go func() { worklist <- []string{pageUrl} }()

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		select {
		case list := <-worklist:
			for _, link := range filterDomain(list) {
				if !seen[link] {
					seen[link] = true
					n++

					log.Debugf("Seen: %s", link)
					go func(link string) {
						worklist <- crawl(link)
					}(link)
				}
				pagesSeen++
			}

		case err := <-httpErrors:
			n++
			log.Println("Fuck! Error from:", err.code, err.url)

		case <-timer.C:
			mux.Lock()
			n++
			log.Println("Pages seen by second: ", pagesSeen - pageSeenLast, ", all: ", pagesSeen)
			pageSeenLast = pagesSeen
			mux.Unlock()
		}
	}
}