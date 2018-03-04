package main

import (
	"strings"
	"net/url"
)

func filterQueue(in chan string, out chan string, done chan struct{}) {
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	var seen = make(map[string]bool)
	for val := range in {
		if !seen[val] {
			seen[val] = true
			wg.Add(1)
			out <- val
		}
	}
}

func validUrl(href string) (bool) {
	return strings.HasPrefix(href, "http")
}

func fixURL(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseURL.ResolveReference(uri)
	return uri.String()
}
