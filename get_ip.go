package main

import (
	"log"
	"net"
	"regexp"
	"strings"
)

var IPv6Regexp = regexp.MustCompile(`\[(\S+)\]`)

type externalIPAddress struct {
	v4 string
	v6 string
}

func isIPValid(addr string) bool {
	if addr != "" {
		ip := net.ParseIP(addr)
		return !(ip == nil)
	}

	return false
}

func getIP(url string, regexp *regexp.Regexp) string {
	var addr string
	body := getURL(url)

	if regexp != nil {
		result := regexp.FindAllStringSubmatch(string(body), -1)
		if len(result) > 0 && len(result[0]) > 0 {
			addr = result[0][1]
		}
	} else {
		addr = strings.Trim(string(body), " \r\n")
	}

	if !isIPValid(addr) {
		addr = ""
	}

	return addr
}

func getExternalIP(conf *config) *externalIPAddress {
	var IPv4, IPv6 string

	IPv4 = getIP("http://ipv4.myexternalip.com/raw", nil)
	if conf.SetIPv6 {
		IPv6 = getIP("http://ipv6.myexternalip.com/raw", nil)
	}

	if IPv4 == "" && IPv6 == "" {
		log.Fatal("couldn't determine external address")
	}

	return &externalIPAddress{v4: IPv4, v6: IPv6}
}
