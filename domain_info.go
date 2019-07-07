package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type domainInfo struct {
	Domain  string `json:"domain"`
	Records []struct {
		RecordID  uint64      `json:"record_id"`
		Type      string      `json:"type"`
		Domain    string      `json:"domain"`
		Fqdn      string      `json:"fqdn"`
		TTL       uint64      `json:"ttl"`
		Subdomain string      `json:"subdomain"`
		Content   string      `json:"content"`
		Priority  interface{} `json:"priority"`
	} `json:"records"`
	Success string `json:"success"`
	Error   string `json:"error"`
}

const (
	getDomainInfoURLTemplate = "https://pddimp.yandex.ru/api2/admin/dns/list?domain=%s"
)

func parseDomainInfoData(data []byte) *domainInfo {
	info := &domainInfo{}

	err := json.Unmarshal(data, info)
	if err != nil {
		log.Fatalf("failed to parse response from Yandex DNS API service %v\n", err)
	}

	return info
}

func getDomainInfo(conf *config) *domainInfo {
	url := fmt.Sprintf(getDomainInfoURLTemplate, conf.Domain)
	body, err := getURL(url, &conf.Token)
	if err != nil {
		log.Fatalf("failed to query '%s': %s", url, err.Error())
	}

	info := parseDomainInfoData(body)
	if info.Success == "error" {
		log.Fatalf("invalid status response: %s\n", info.Error)
	}

	return info
}
