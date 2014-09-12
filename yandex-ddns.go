package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var (
		configFile     string
		testConfigOnly bool
	)

	flag.StringVar(&configFile, "config", "yandex-ddns.json", "configuration file")
	flag.BoolVar(&testConfigOnly, "t", false, "only test configuration file")

	flag.Parse()

	conf := newConfigurationFromFile(configFile)

	if testConfigOnly {
		verifyConfiguration(conf)
		fmt.Println("Configuration file Ok.")
		os.Exit(0)
	}

	if conf.LogFile != "" {
		f, err := os.OpenFile(conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	verifyConfiguration(conf)
	extIPAddr := getExternalIP()
	domainInfo := getDomainInfo(conf)
	updateDomainAddress(domainInfo, extIPAddr, conf)
}
