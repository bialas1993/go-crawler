package main

import (
	log "github.com/sirupsen/logrus"
		"github.com/bialas1993/go-crawler/crawler"
)

func main() {
	pageUrl, depth, logTimer, logLevel := parseParams()

	log.Printf("pageUrl=%s, depth=%d, logTimer=%d, logLevel=%d", pageUrl, depth, logTimer, logLevel)

	crawler.Run(pageUrl)
}


//func createSeoFilters() FilterManager {
//	fm := FilterManager{}
//
//	fm.
//		Add(filters.HTag{}).
//		Add(filters.ImgTag{})
//
//	return fm
//}