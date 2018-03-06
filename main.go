package main

import (
	"time"
	"sync"
)

const CONNECTIONS_LIMIT = 25

var tokens = make(chan struct{}, CONNECTIONS_LIMIT)
var httpErrors = make(chan *HttpError)
var mux sync.Mutex

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
	page, _, logTimer := parseParams()
	worklist := make(chan []string)
	var n int
	var pagesSeen, pageSeenLast uint16
	var timer *time.Ticker

	if logTimer > 0 {
		timer = time.NewTicker( time.Duration(logTimer) * time.Second )
	} else {
		timer = time.NewTicker( time.Second )
		timer.Stop()
	}

	n++
	go func() { worklist <- []string{page} }()

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		select {
		case list := <-worklist:
			for _, link := range list {
				if !seen[link] {
					seen[link] = true
					n++
					go func(link string) {
						worklist <- crawl(link)
					}(link)
				}
				pagesSeen++
			}

		case err := <-httpErrors:
			n++
			println("Fuck! Error from:", err.code, err.url)

		case <-timer.C:
			mux.Lock()
			n++
			println("Pages seen by second: ", pagesSeen - pageSeenLast, ", all: ", pagesSeen)
			pageSeenLast = pagesSeen
			mux.Unlock()
		}
	}
}