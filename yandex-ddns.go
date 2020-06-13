package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nightlyone/lockfile"
)

func initLock(file string) (*lockfile.Lockfile, error) {
	lock, err := lockfile.New(filepath.Join(os.TempDir(), file))
	if err != nil {
		return nil, err
	}

	err = lock.TryLock()
	if err != nil {
		return nil, err
	}

	return &lock, nil
}

func main() {
	var (
		configFile     string
		testConfigOnly bool
	)

	flag.StringVar(&configFile, "config", "yandex-ddns.toml", "configuration file")
	flag.BoolVar(&testConfigOnly, "t", false, "only test configuration file")

	flag.Parse()

	conf := newConfigurationFromFile(configFile)

	if testConfigOnly {
		verifyConfiguration(conf)
		fmt.Println("Configuration file Ok.")
		os.Exit(0)
	}

	lock, err := initLock("yandex-ddns.lock")
	if err != nil {
		log.Fatalf("Couldn't init lock file: %v\n", err)
	}
	defer lock.Unlock()

	if conf.LogFile != "" {
		f, err := os.OpenFile(conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer closeResource(f)

		log.SetOutput(f)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	verifyConfiguration(conf)
	extIPAddr := getExternalIP(conf)
	domainInfo := getDomainInfo(conf)
	updateDomainAddress(domainInfo, extIPAddr, conf)
}
