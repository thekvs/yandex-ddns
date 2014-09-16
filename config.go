package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type config struct {
	Token     string  `json:"token"`
	Domain    string  `json:"domain"`
	SubDomain string  `json:"subdomain"`
	LogFile   string  `json:"logfile"`
	TTL       *uint64 `json:"ttl,omitempty"`
	SetIPv6   bool    `json:"set-ipv6"`
}

const (
	minTTLValue = 900
	maxTTLValue = 1209600
)

func newConfigurationFromFile(path string) *config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("can't open configuration file: %v\n", err)
	}

	conf := newConfiguration(file)

	return conf
}

func verifyConfiguration(conf *config) {
	if conf.Token == "" {
		log.Fatal("missed mandatory configuration parameter 'token'")
	}

	if conf.Domain == "" {
		log.Fatal("missed mandatory configuration parameter 'domain'")
	}

	if conf.TTL != nil {
		if *conf.TTL < minTTLValue || *conf.TTL > maxTTLValue {
			log.Fatalf("domain TTL value (=%d) exeeds permissible range (=[%d, %d])\n", *conf.TTL, minTTLValue, maxTTLValue)
		}
	}
}

func newConfiguration(data io.Reader) *config {
	decoder := json.NewDecoder(data)
	conf := &config{}
	decoder.Decode(&conf)

	return conf
}
