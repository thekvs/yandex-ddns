package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	userAgentHeader  = "User-Agent"
	defaultUserAgent = "Mozilla/4.0 (compatible; MSIE 7.0; +https://github.com/thekvs/yandex-ddns)"
)

const defaultNetworkTimeout = 20 * time.Second

var client = &http.Client{Timeout: defaultNetworkTimeout}

func getURL(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("error: '%v'\n", err)
	}

	req.Header.Add(userAgentHeader, defaultUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error: '%v'\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error: '%v'\n", err)
	}

	return body
}
