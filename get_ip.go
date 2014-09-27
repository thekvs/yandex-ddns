package main

import (
	"log"
	"net"
	"regexp"
)

var IPv4Regexp = regexp.MustCompile("IPv4: (\\S+)")
var IPv6Regexp = regexp.MustCompile("IPv6: (\\S+)")

type externalIPAddress struct {
	v4 string
	v6 string
}

func isIPValid(addr string) bool {
	if addr != "" {
		ip := net.ParseIP(addr)
		if ip == nil {
			return false
		}
		return true
	}

	return false
}

func getIP(url string, regexp *regexp.Regexp) string {
	body := getURL(url)
	result := regexp.FindAllStringSubmatch(string(body), -1)

	var addr string
	if len(result) > 0 && len(result[0]) > 0 {
		addr = result[0][1]
		if !isIPValid(addr) {
			addr = ""
		}
	}

	return addr
}

func getExternalIP(conf *config) *externalIPAddress {
	var IPv4, IPv6 string

	IPv4 = getIP("https://ipv4.internet.yandex.ru/", IPv4Regexp)
	if conf.SetIPv6 {
		IPv6 = getIP("https://ipv6.internet.yandex.ru/", IPv6Regexp)
	}

	if IPv4 == "" && IPv6 == "" {
		log.Fatal("couldn't determine external address")
	}

	return &externalIPAddress{v4: IPv4, v6: IPv6}
}
