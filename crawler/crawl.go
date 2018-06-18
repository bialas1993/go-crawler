package crawler

import (
	"net/url"
	log "github.com/sirupsen/logrus"
)

const (
	CONNECTIONS_LIMIT = 20
)


var n uint32

type crawler struct{
	tokens     chan struct{}
	httpErrors chan *HttpError
	page       *url.URL
	workList   chan []string
}

func Run(pageUrl string) *crawler {
	var page, _ = url.Parse(pageUrl)

	c := crawler{
		tokens:     make(chan struct{}, CONNECTIONS_LIMIT),
		httpErrors: make(chan *HttpError),
		page:       page,
		workList:   make(chan []string),
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

						log.Debugf("Seen: %s", link)
						go func(link string) {
							c.workList <- c.crawl(link)
						}(link)

						log.Printf("Seen: %s", link)
					}
				}
			} else {
				break
			}


		case err := <-c.httpErrors:
			n++
			log.Println("Fuck! Error from:", err.code, err.url)
		}
	}

	defer func(c *crawler) {
		close(c.workList)
		close(c.tokens)
		close(c.httpErrors)
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
