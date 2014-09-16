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
	addr := regexp.FindAllStringSubmatch(string(body), -1)[0][1]
	if !isIPValid(addr) {
		addr = ""
	}

	return addr
}

func getExternalIP() *externalIPAddress {
	IPv4 := getIP("http://ipv4.internet.yandex.ru/", IPv4Regexp)
	IPv6 := getIP("http://ipv6.internet.yandex.ru/", IPv6Regexp)

	if IPv4 == "" && IPv6 == "" {
		log.Fatal("Coudn't get neither IPv4 nor IPv6 addresses")
	}

	return &externalIPAddress{v4: IPv4, v6: IPv6}
}
