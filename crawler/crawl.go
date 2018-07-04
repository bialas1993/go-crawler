package crawler

import (
	"net/url"
	"golang.org/x/net/html"
	"context"
	"net/http"
	"encoding/base64"
)

const (
	CONNECTIONS_LIMIT = 200
	AUTH_DISABLE      = 0
	AUTH_ENABLE       = 1
)

type AuthCredentials struct {
	Enabled  bool
	Username string
	Password string
}

func (a AuthCredentials) Hash() string {
	return base64.StdEncoding.EncodeToString([]byte(a.Username + ":" + a.Password))
}

type crawlUrl struct {
	Parent string
	Url    string
}

type crawler struct {
	tokens      chan struct{}
	httpErrors  chan *HttpError
	workList    chan []crawlUrl
	page        *url.URL
	nodeChannel chan *html.Node
	ctx         context.Context
	cancel      context.CancelFunc
	logger      LogService
	auth        AuthCredentials
}

var n uint32 = 1

func Run(pageUrl string, nodeChannel *chan *html.Node, ctx context.Context, cancel context.CancelFunc, logger LogService, auth AuthCredentials) {
	var page, _ = url.Parse(pageUrl)

	c := crawler{
		tokens:      make(chan struct{}, CONNECTIONS_LIMIT),
		httpErrors:  make(chan *HttpError),
		workList:    make(chan []crawlUrl),
		page:        page,
		nodeChannel: *nodeChannel,
		ctx:         ctx,
		cancel:      cancel,
		logger:      logger,
		auth:        auth,
	}

	go func() {
		c.workList <- []crawlUrl{crawlUrl{
			Parent: "",
			Url: pageUrl,
		}}
	}()

	c.bind()
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
							c.logger.Log(Debug("Seen: "+link.Url, ExtraDataMessage{Url: link.Url, ResponseCode: http.StatusOK}))
						}(link)
					}
				}
				continue
			}
			break

		case err := <-c.httpErrors:
			n++
			go func(err *HttpError) {
				c.logger.Log(Error(err.Error(), ExtraDataMessage{Url: err.url, ResponseCode: err.code}))
			}(err)
		}
	}

	c.logger.Log(Info("Not found any more pages to see"))
	c.cancel()

	go func() {
		c.nodeChannel <- nil
	}()

	defer func(c *crawler) {
		close(c.workList)
		close(c.tokens)
		close(c.httpErrors)
		c.logger.Close()
	}(c)
}

func (c *crawler) crawl(url crawlUrl) []crawlUrl {
	c.tokens <- struct{}{}
	list, err := extract(c.nodeChannel, url, c.auth)
	<-c.tokens

	if err != nil {
		err, ok := err.(*HttpError)
		if ok {
			go func(err *HttpError) {
				c.httpErrors <- err
			}(err)
		}
	}

	return list
}
