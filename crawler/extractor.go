package crawler

import (
	"net/http"
	"crypto/tls"
	"golang.org/x/net/html"
	"fmt"
	"regexp"
	"strings"
	"net/url"
)

func extract(nodeChannel chan *html.Node, url crawlUrl, auth AuthCredentials) ([]crawlUrl, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url.Url, nil)

	if auth.Enabled {
		req.Header.Set("Authorization", "Basic " + auth.Hash())
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, &HttpError{url.Parent, url.Url, resp.StatusCode}
	}

	doc, err := html.Parse(resp.Body)
	go func() {
		nodeChannel <- doc
	}()

	if err != nil {
		return nil, fmt.Errorf("parsing %s to HTML: %v", url, err)
	}

	var links []crawlUrl
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue
				}
				links = append(links, crawlUrl{
					Parent: url.Url,
					Url: link.String(),
				})
			}
		}
	}

	forEachNode(doc, visitNode, nil)

	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func filterDomain(page *url.URL, list []crawlUrl) []crawlUrl {
	var re = regexp.MustCompile(`(?m)^(http(|s):|)(\/\/)((.*@)|)` + page.Hostname() + `.*`)
	var splitedUrl []string
	var urls []crawlUrl

	for _, url := range list {
		splitedUrl = strings.Split(url.Url, "#")
		url.Url = splitedUrl[0]

		if len(re.FindAllString(url.Url, -1)) > 0 {
			urls = append(urls, url)
		}
	}

	return urls
}

