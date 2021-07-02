package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const stringToSearch = "concurrency"

var sites = []string{
	"https://google.com",
	"https://itc.ua/",
	"https://twitter.com/concurrencyinc",
	"https://twitter.com/",
	"http://localhost:8000",
	"https://github.com/bradtraversy/go_restapi/blob/master/main.go",
	"https://www.youtube.com/",
	"https://postman-echo.com/get",
	"https://en.wikipedia.org/wiki/Concurrency_(computer_science)#:~:text=In%20computer%20science%2C%20concurrency%20is,without%20affecting%20the%20final%20outcome.",
}

type SiteData struct {
	data []byte
	uri  string
}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	resultsCh := make(chan SiteData, len(sites))

	go worker(ctx, resultsCh, sites)

	go reader(cancel, resultsCh)
	// give one second to validate if all other goroutines are closed
	time.Sleep(time.Second)
}

func worker(ctx context.Context, siteData chan SiteData, sites []string) {
	for _, u := range sites {
		go makeRequest(ctx, siteData, u)
	}
}

func reader(cancel context.CancelFunc, siteDataChan chan SiteData) {
	defer fmt.Println("exiting from searcher...")
	for data := range siteDataChan {
		if strings.Contains(string(data.data), stringToSearch) {
			fmt.Printf("'%s' string is found in %s \n", stringToSearch, data.uri)
			cancel()
			return
		} else {
			fmt.Printf("Nothing found in  %s \n", data.uri)
		}
	}
}

func makeRequest(ctx context.Context, siteData chan SiteData, uri string) {
	fmt.Printf("starting sending request to %s \n", uri)
	select {
	case <-ctx.Done():
		fmt.Printf("Get %s : context canceled\n", uri)

	default:
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		siteData <- SiteData{
			data: bodyBytes,
			uri:  uri,
		}
	}
}
