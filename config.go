package main

import (
	"io"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	Token     string  `toml:"token"`
	Domain    string  `toml:"domain"`
	SubDomain string  `toml:"subdomain"`
	LogFile   string  `toml:"logfile"`
	TTL       *uint64 `toml:"ttl,omitempty"`
	SetIPv6   bool    `toml:"set-ipv6"`
}

const (
	minTTLValue = 900
	maxTTLValue = 1209600
)

const allowedPermissions = 0600

func isPermissionsOk(f *os.File) bool {
	finfo, err := f.Stat()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	permissions := finfo.Mode().Perm()
	if permissions == allowedPermissions {
		return true
	}

	return false
}

func newConfigurationFromFile(path string) *config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("can't open configuration file: %v\n", err)
	}
	defer file.Close()

	if !isPermissionsOk(file) {
		log.Fatalf("error: configuration file with sensitive information has insecure permissions\n")
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
	var conf config
	if _, err := toml.DecodeReader(data, &conf); err != nil {
		log.Fatalf("Couldn't parse configuration file: %v", err)
	}

	return &conf
}
