package main

import (
	"log"
	"net"
	"regexp"
	"strings"
)

type lookupExternalIPUrl struct {
	v4 string
	v6 string
}

type externalIPAddress struct {
	v4 string
	v6 string
}

var lookupExternalIPUrls = []lookupExternalIPUrl{
	lookupExternalIPUrl{
		v4: "http://ipv4.myexternalip.com/raw",
		v6: "http://ipv6.myexternalip.com/raw",
	},
	lookupExternalIPUrl{
		v4: "https://v4.ifconfig.co/ip",
		v6: "https://v6.ifconfig.co/ip",
	},
}

func isIPValid(addr string) bool {
	if addr != "" {
		ip := net.ParseIP(addr)
		return !(ip == nil)
	}

	return false
}

func getIP(url string, regexp *regexp.Regexp) (string, error) {
	var addr string

	body, err := getURL(url)
	if err != nil {
		return "", err
	}

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

	return addr, nil
}

func getExternalIP(conf *config) *externalIPAddress {
	var IPv4, IPv6 string
	var err error

	for _, lookup := range lookupExternalIPUrls {
		IPv4, err = getIP(lookup.v4, nil)
		if err != nil {
			log.Printf("%s", err.Error())
		}

		if conf.SetIPv6 {
			IPv6, err = getIP(lookup.v6, nil)
			if err != nil {
				log.Printf("%s", err.Error())
			}
		}

		if len(IPv4) > 0 || len(IPv6) > 0 {
			break
		}
	}

	if len(IPv4) == 0 && len(IPv6) == 0 {
		log.Fatal("couldn't determine external address")
	}

	return &externalIPAddress{v4: IPv4, v6: IPv6}
}
