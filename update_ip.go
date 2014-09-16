package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
)

const (
	editARecordURLTemplate    = "https://pddimp.yandex.ru/nsapi/edit_a_record.xml?token=%s&domain=%s&subdomain=%s&record_id=%s&ttl=%d&content=%s"
	editAAAARecordURLTemplate = "https://pddimp.yandex.ru/nsapi/edit_aaaa_record.xml?token=%s&domain=%s&subdomain=%s&record_id=%s&ttl=%d&content=%s"
)

type updateRecordResponse struct {
	Name  string `xml:"domains>domain>name"`
	Error string `xml:"domains>error"`
}

func verifyUpdateRecordResponse(data []byte) {
	resp := &updateRecordResponse{}

	err := xml.Unmarshal(data, resp)
	if err != nil {
		log.Fatalf("failed to parse response from Yandex DNS API service %v\n", err)
	}

	if resp.Error != "ok" {
		log.Fatalf("update failed, error message: %v\n", resp.Error)
	}
}

func getFullDomainName(subdomain string, domain string) string {
	if subdomain == "@" {
		return domain
	}

	return fmt.Sprintf("%s.%s", subdomain, domain)
}

func updateDomainAddress(info *domainInfo, extIPAddr *externalIPAddress, conf *config) {
	subDomain := conf.SubDomain
	if conf.SubDomain == "" {
		subDomain = "@"
	}

	var (
		url  string
		addr string
	)

	update := func() {
		body := getURL(url)
		verifyUpdateRecordResponse(body)

		log.Printf("IP address for '%s' set to %s\n", getFullDomainName(subDomain, conf.Domain), addr)
	}

	var ttl uint64

	updated := false
	for _, record := range info.Records {
		if record.SubDomain != subDomain {
			continue
		}

		if conf.TTL != nil && *conf.TTL > 0 {
			ttl = *conf.TTL
		} else {
			var err error
			ttl, err = strconv.ParseUint(record.TTL, 10, 0)
			if err != nil {
				log.Fatalf("failed to convert %s to uint: %v\n", record.TTL, err)
			}
		}

		if record.Type == "A" && extIPAddr.v4 != "" {
			addr = extIPAddr.v4
			url = fmt.Sprintf(editARecordURLTemplate, conf.Token, conf.Domain, record.SubDomain, record.ID, ttl, addr)
		} else if record.Type == "AAAA" && extIPAddr.v6 != "" {
			addr = extIPAddr.v6
			url = fmt.Sprintf(editAAAARecordURLTemplate, conf.Token, conf.Domain, record.SubDomain, record.ID, ttl, addr)
		} else {
			continue
		}

		update()
		updated = true
	}

	if !updated {
		log.Fatalf("domain '%s' not known for Yandex.DNS\n", getFullDomainName(subDomain, conf.Domain))
	}
}
