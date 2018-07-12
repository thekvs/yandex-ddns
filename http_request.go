package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	userAgentHeader  = "User-Agent"
	defaultUserAgent = "Mozilla/4.0 (compatible; MSIE 7.0; +https://github.com/thekvs/yandex-ddns)"
)

const defaultNetworkTimeout = 20 * time.Second

var client = &http.Client{Timeout: defaultNetworkTimeout}

func getURL(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(userAgentHeader, defaultUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeResource(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpecetd HTTP status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
