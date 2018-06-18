package crawler

import (
	"net/url"
		"strconv"
)

const (
	CONNECTIONS_LIMIT = 20
	CLOSE_LOGGER_MESSAGE = "crawler_logger_exit;"
)

var n uint32

type crawler struct{
	tokens     chan struct{}
	httpErrors chan *HttpError
	workList   chan []string
	page       *url.URL
	logChannel chan string
}

func Run(pageUrl string, logChannel *chan string) *crawler {
	var page, _ = url.Parse(pageUrl)

	c := crawler{
		tokens:     make(chan struct{}, CONNECTIONS_LIMIT),
		httpErrors: make(chan *HttpError),
		workList:   make(chan []string),
		page:       page,
		logChannel: *logChannel,
	}

	go func() {
		c.workList <- []string{pageUrl}
	}()

	n++
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
					if !seen[link] {
						seen[link] = true
						n++

						go func(link string) {
							c.workList <- c.crawl(link)
							//c.logChannel <- "Seen: " +  link
						}(link)
					}
				}
			} else {
				break
			}
		case err := <-c.httpErrors:
			n++
			go func(err *HttpError) {
				c.logChannel <- "Fuck! Error from: " + strconv.Itoa(err.code) + " " + err.url
			} (err)
		}
	}

	defer func(c *crawler) {
		close(c.workList)
		close(c.tokens)
		close(c.httpErrors)
		c.logChannel <- CLOSE_LOGGER_MESSAGE
	}(c)
}

func (c *crawler) crawl(url string) []string {
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
