package crawler

import (
	"net/url"
	"golang.org/x/net/html"
	"context"
)

const (
	CONNECTIONS_LIMIT = 200
)

var n uint32 = 1

type crawlUrl struct{
	Parent string
	Url string
}

type crawler struct{
	tokens      chan struct{}
	httpErrors  chan *HttpError
	workList    chan []crawlUrl
	page        *url.URL
	logChannel  chan LogMessage
	nodeChannel chan *html.Node
	ctx         context.Context
	cancel      context.CancelFunc
}

func Run(pageUrl string, nodeChannel *chan *html.Node, logChannel *chan LogMessage, ctx context.Context, cancel context.CancelFunc) *crawler {
	var page, _ = url.Parse(pageUrl)

	c := crawler{
		tokens:      make(chan struct{}, CONNECTIONS_LIMIT),
		httpErrors:  make(chan *HttpError),
		workList:    make(chan []crawlUrl),
		page:        page,
		logChannel:  *logChannel,
		nodeChannel: *nodeChannel,
		ctx:         ctx,
		cancel:      cancel,
	}

	go func() {
		c.workList <- []crawlUrl{crawlUrl{
			Parent: "",
			Url: pageUrl,
		}}
	}()

	c.bind()

	return &c
}

func (c *crawler) bind() {
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		select {
		case list := <-c.workList:
			if len(list) > 0 {
				for _, link := range filterDomain(c.page, list) {
					if !seen[link.Url] {
						seen[link.Url] = true
						n++

						go func(link crawlUrl) {
							c.workList <- c.crawl(link)
							c.logChannel <- Debug("Seen: " + link.Url)
						}(link)
					}
				}
				continue
			}
			break

		case err := <-c.httpErrors:
			n++
			go func(err *HttpError) {
				c.logChannel <- Error(err.Error())
			} (err)
		}
	}

	c.logChannel <- Info("Not found any more pages to see")
	c.cancel()

	go func() {
		c.logChannel <- CloseLogger()
		c.nodeChannel <- nil
	}()

	defer func(c *crawler) {
		close(c.workList)
		close(c.tokens)
		close(c.httpErrors)
	}(c)
}

func (c *crawler) crawl(url crawlUrl) []crawlUrl {
	c.tokens <- struct{}{}
	list, err := extract(c.nodeChannel, url)
	<-c.tokens

	if err != nil {
		err, ok := err.(*HttpError); if ok {
			go func (err *HttpError) {
				c.httpErrors <- err
			}(err)
		}
	}

	return list
}
