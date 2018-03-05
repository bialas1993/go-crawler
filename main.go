package main

import (
	"fmt"
	"log"
)

var tokens = make(chan struct{}, 1)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{}
	list, err := Extract(url)
	<-tokens

	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	page, depth := parseParams()
	worklist := make(chan []string)
	var n int
	lvl := 0

	n++
	go func() { worklist <- []string{page} }()

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist

		for _, link := range list {
			if lvl >= depth {
				return
			}

			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}

		lvl++
	}
}