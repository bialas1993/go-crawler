package main

import (
	"fmt"
	"sync"
	"crypto/tls"
	"github.com/jackdanger/collectlinks"
	"net/http"
	"bytes"
	"time"
	"net"
)

type writer struct {
	fileName string
	body     []byte
}

var writes = make(chan writer)
//var wg BoundedWaitGroup
var wg sync.WaitGroup
var seens = 0

func main() {
	args := validateParams()
	done := make(chan struct{})
	queue := make(chan string, 3)
	filteredQueue := make(chan string)
	//wg = NewBoundedWaitGroup(2)

	go func() { queue <- args[0] }()
	go filterQueue(queue, filteredQueue, done)

	go func() {
		for uri := range filteredQueue {
			go enqueue(uri, queue)
		}
	}()

	for {
		select {
		case data := <-writes:
			seens++
			println("Seens: ", seens)
			fmt.Println("write file:", data.fileName)
		case <-done:
			fmt.Println("crawl done")
			return
		}
	}
}

func enqueue(uri string, queue chan string) {
	fmt.Println("fetching", uri)

	if validUrl(uri) {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Dial: TimeoutDialer(&Config{
				ConnectTimeout:   30 * time.Second,
				ReadWriteTimeout: 30 * time.Second,
			}),
		}
		client := http.Client{Transport: transport}
		resp, err := client.Get(uri)

		if err != nil {
			println("[ERROR] Someting wen't wrong! ", uri)
			wg.Done()

			return
		}

		if resp.StatusCode >= 400 {
			println("Error fetch ", uri)
			wg.Done()

			return
		}


		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		defer resp.Body.Close()

		responseBytes := buf.Bytes()
		links := collectlinks.All(bytes.NewReader(responseBytes))

		writes <- writer{fileName: uri, body: responseBytes}

		if len(links) == 0 {
			wg.Done()
		}

		for _, link := range links {
			absolute := fixURL(link, uri)
			if uri != "" {
				go func() { queue <- absolute }()
			}
		}
	}
}


type Config struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

func TimeoutDialer(config *Config) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, config.ConnectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(config.ReadWriteTimeout))
		return conn, nil
	}
}

