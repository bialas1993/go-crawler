package crawler

import (
	"net/url"
	)

const (
	CONNECTIONS_LIMIT = 100
	CLOSE_LOGGER_MESSAGE = "crawler_logger_exit;"
)

var n uint32 = 1

type crawlUrl struct{
	Parent string
	Url string
}

type crawler struct{
	tokens        chan struct{}
	httpErrors    chan *HttpError
	workList      chan []crawlUrl
	page          *url.URL
	logChannel    chan string
}

func Run(pageUrl string, logChannel *chan string) *crawler {
	var page, _ = url.Parse(pageUrl)

	c := crawler{
		tokens:        make(chan struct{}, CONNECTIONS_LIMIT),
		httpErrors:    make(chan *HttpError),
		workList:      make(chan []crawlUrl),
		page:          page,
		logChannel:    *logChannel,
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
							c.logChannel <- "Seen: " + link.Url
						}(link)
					}
				}

				continue
			}
			break

		case err := <-c.httpErrors:
			n++
			go func(err *HttpError) {
				c.logChannel <- err.Error()
			} (err)
		}
	}

	c.logChannel <- "Not found any more pages to see"

	defer func(c *crawler) {
		close(c.workList)
		close(c.tokens)
		close(c.httpErrors)
		c.logChannel <- CLOSE_LOGGER_MESSAGE
	}(c)
}

func (c *crawler) crawl(url crawlUrl) []crawlUrl {
	c.tokens <- struct{}{}
	list, err := extract(url)
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
