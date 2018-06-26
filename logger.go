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
			if msg.Level() != crawler.LOG_MESSAGE_CLOSE {
				switch msg.Level() {
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
			}
			break
		}
	}
}

func (LoggerService) Panic(message string, data interface{}) {
	panic("implement me")
}

func (LoggerService) Fatal(message string, data interface{}) {
	panic("implement me")
}

func (LoggerService) Error(message string, data interface{}) {
	panic("implement me")
}

func (LoggerService) Warning(message string, data interface{}) {
	panic("implement me")
}

func (LoggerService) Info(message string, data interface{}) {
	panic("implement me")
}

func (LoggerService) Debug(message string, data interface{}) {
	panic("implement me")
}

func (ls *LoggerService) Log(message crawler.LogMessage) {
	go func() {ls.Stream <- message} ()
}

func (ls *LoggerService) Close() {
	close(ls.Stream)
}


