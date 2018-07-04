package main

import (
	"github.com/bialas1993/go-crawler/crawler"
	"log"
	"time"
	"context"
	"gopkg.in/olivere/elastic.v6"
	"github.com/google/uuid"
	"os"
)

const (
	ELASTICSEARCH_LOGGER_HOST_ENV = "ELASTICSEARCH_LOGGER_HOST"
	ELASTICSEARCH_LOGGER_INDEX = "crawler"
	ELASTICSEARCH_LOGGER_MAPPING = `{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"crawler":{
			"properties":{
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"created":{
					"type":"date"
				},
				"code":{
					"type":"integer"
				},
				"domain":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"response_code":{
					"type":"integer"
				}
			}
		}
	}
}`
)

type ElasticMessage struct {
	UUID		 string    `json:"uuid"`
	Domain 		 string    `json:"domain"`
	Message      string    `json:"message"`
	Code         int       `json:"code"`
	Created      time.Time `json:"created,omitempty"`
	Url          string    `json:"url,omitempty"`
	ResponseCode int       `json:"response_code,omitempty"`
}

func CreateElasticLoggerService(globalMsg crawler.ExtraDataMessage) *ElasticLoggerService {
	l := new(ElasticLoggerService)
	l.Stream = make(chan crawler.LogMessage)
	l.globalMessage = globalMsg
	l.UUID = uuid.New().String()

	l.configure()

	go func() { l.Listen() } ()

	return l
}

func (l *ElasticLoggerService) configure() {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(os.Getenv(ELASTICSEARCH_LOGGER_HOST_ENV)))
	if err != nil {
		panic(err)
	}

	client.Ping(os.Getenv(ELASTICSEARCH_LOGGER_HOST_ENV)).Do(ctx)
	if err != nil {
		panic(err)
	}

	l.client = client
	exists, err := client.IndexExists(ELASTICSEARCH_LOGGER_INDEX).Do(ctx)

	if err != nil {
		log.Println(err.Error())
	}

	if !exists {
		createIndex, err := client.
				CreateIndex(ELASTICSEARCH_LOGGER_INDEX).
				BodyString(ELASTICSEARCH_LOGGER_MAPPING).
				Do(ctx)

		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

type ElasticLoggerService struct {
	LoggerService
	client *elastic.Client
	globalMessage crawler.ExtraDataMessage
	UUID string
}

func (l ElasticLoggerService) Log(message crawler.LogMessage) {
	msg := ElasticMessage{
		Message: message.GetMessage(),
		Code: message.GetLevel(),
		Created: time.Now(),
	}
	
	data := message.GetExtraData()
	msg.Url = data.Url
	msg.ResponseCode = data.ResponseCode
	msg.Domain = l.globalMessage.Domain
	msg.UUID = l.UUID

	l.client.Index().
		Index(ELASTICSEARCH_LOGGER_INDEX).
		Type("crawler").
		BodyJson(msg).
		Do(context.Background())

	l.LoggerService.Log(message)
}
