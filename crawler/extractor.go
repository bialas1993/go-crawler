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

func extract(url string) ([]string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, &HttpError{url, resp.StatusCode}
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing %s to HTML: %v", url, err)
	}

	var links []string
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
				links = append(links, link.String())
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

func filterDomain(page *url.URL, list []string) []string {
	var re = regexp.MustCompile(`(?m)^(http(|s):|)(\/\/)` + page.Hostname() + `.*`)
	var splitedUrl []string
	var urls []string

	for _, url := range list {
		splitedUrl = strings.Split(url, "#")
		url = splitedUrl[0]

		if len(re.FindAllString(url, -1)) > 0 {
			urls = append(urls, url)
		}
	}

	return urls
}

