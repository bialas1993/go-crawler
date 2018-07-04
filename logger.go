package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/bialas1993/go-crawler/crawler"
)

type LoggerService struct{
	Stream chan crawler.LogMessage
}

func CreateLogger() *LoggerService {
	ls := new(LoggerService)
	ls.Stream = make(chan crawler.LogMessage)

	go func() {
		ls.Listen()
	} ()

	return ls
}

func (ls *LoggerService) Listen() {
	for {
		select {
		case msg := <-ls.Stream:
			//if msg.GetLevel() != crawler.LOG_MESSAGE_CLOSE {
				switch msg.GetLevel() {
				case crawler.LOG_MESSAGE_INFO:
					log.Info(msg.GetMessage())
				case crawler.LOG_MESSAGE_WARN:
					log.Warn(msg.GetMessage())
				case crawler.LOG_MESSAGE_ERROR:
					log.Error(msg.GetMessage())
				case crawler.LOG_MESSAGE_DEBUG:
					log.Debug(msg.GetMessage())
				}
				continue
			//}
			//break
		}
	}
}

func (ls *LoggerService) Log(message crawler.LogMessage) {
	go func() {ls.Stream <- message} ()
}

func (ls *LoggerService) Close() {
	close(ls.Stream)
}


