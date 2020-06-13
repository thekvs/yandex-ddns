package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

const (
	editRecordURL = "https://pddimp.yandex.ru/api2/admin/dns/edit"
)

type updateDomainResponse struct {
	Domain   string `json:"domain"`
	RecordID uint64 `json:"record_id"`
	Record   struct {
		RecordID  uint64      `json:"record_id"`
		Type      string      `json:"type"`
		Domain    string      `json:"domain"`
		Fqdn      string      `json:"fqdn"`
		TTL       uint64      `json:"ttl"`
		Subdomain string      `json:"subdomain"`
		Content   string      `json:"content"`
		Priority  interface{} `json:"priority"`
		Operation string      `json:"operation"`
	} `json:"record"`
	Success string `json:"success"`
	Error   string `json:"error"`
}

func verifyUpdateRecordResponse(data []byte) {
	resp := &updateDomainResponse{}

	err := json.Unmarshal(data, resp)
	if err != nil {
		log.Fatalf("failed to parse response from Yandex DNS API service: %s\n", err.Error())
	}

	if resp.Success != "ok" {
		log.Fatalf("update failed: %s\n", resp.Error)
	}
}

func getFullDomainName(subdomain, domain string) string {
	if subdomain == "@" {
		return domain
	}

	return fmt.Sprintf("%s.%s", subdomain, domain)
}

func update(token string, params *url.Values) bool {
	resp, err := postURL(editRecordURL, &token, params)
	if err != nil {
		log.Fatalf("%s", err.Error())
	} else {
		verifyUpdateRecordResponse(resp)
	}

	return true
}

func updateDomainAddress(info *domainInfo, extIPAddr *externalIPAddress, conf *config) {
	subdomain := conf.SubDomain
	if len(conf.SubDomain) == 0 {
		subdomain = "@"
	}

	var (
		ttl  string
		addr string
	)

	updated := false

	for _, record := range info.Records {
		if record.Subdomain != subdomain {
			continue
		}

		if conf.TTL != nil && *conf.TTL > 0 {
			ttl = fmt.Sprintf("%d", *conf.TTL)
		} else {
			ttl = fmt.Sprintf("%d", record.TTL)
		}

		if record.Type == "A" && len(extIPAddr.v4) > 0 {
			addr = extIPAddr.v4
		} else if record.Type == "AAAA" && len(extIPAddr.v6) > 0 {
			addr = extIPAddr.v6
		} else {
			continue
		}

		params := &url.Values{}

		params.Set("domain", conf.Domain)
		params.Set("record_id", fmt.Sprintf("%d", record.RecordID))
		params.Set("ttl", ttl)
		params.Set("content", addr)

		updated = update(conf.Token, params)

		log.Printf("IP address for '%s' set to '%s'\n", conf.Domain, addr)
	}

	if !updated {
		log.Fatalf("domain '%s' not known to Yandex.DNS\n", getFullDomainName(subdomain, conf.Domain))
	}
}
