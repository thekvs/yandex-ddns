package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	userAgentHeader  = "User-Agent"
	defaultUserAgent = "Mozilla/4.0 (compatible; MSIE 7.0; +https://github.com/thekvs/yandex-ddns)"
)

const defaultNetworkTimeout = 20 * time.Second

var client = &http.Client{Timeout: defaultNetworkTimeout}

func postURL(url string, token *string, values *url.Values) ([]byte, error) {
	var (
		req *http.Request
		err error
	)

	if values != nil {
		req, err = http.NewRequest("POST", url, bytes.NewBufferString(values.Encode()))
	} else {
		req, err = http.NewRequest("POST", url, nil)
	}
	if err != nil {
		log.Fatalf("%s", err)
	}

	if token != nil && len(*token) > 0 {
		req.Header.Add("pddToken", *token)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("%s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("%s", err)
	}

	return body, nil
}

func getURL(url string, token *string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(userAgentHeader, defaultUserAgent)

	if token != nil && len(*token) > 0 {
		req.Header.Add("pddToken", *token)
	}

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
