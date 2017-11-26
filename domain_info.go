package main

import (
	"encoding/xml"
	"fmt"
	//	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type domainInfoRecord struct {
	Domain    string `xml:"domain,attr"`
	Priority  string `xml:"priority,attr"`
	TTL       string `xml:"ttl,attr"`
	SubDomain string `xml:"subdomain,attr"`
	Type      string `xml:"type,attr"`
	ID        string `xml:"id,attr"`
}

type domainInfo struct {
	Name      string             `xml:"domains>domain>name"`
	Records   []domainInfoRecord `xml:"domains>domain>response>record"`
	Delegated *string            `xml:"domains>domain>nsdelegated,omitempty"`
	Error     string             `xml:"domains>error"`
}

const (
	getDomainInfoURLTemplate = "https://pddimp.yandex.ru/nsapi/get_domain_records.xml?token=%s&domain=%s"
)

func parseDomainInfoData(data []byte) *domainInfo {
	info := &domainInfo{}

	err := xml.Unmarshal(data, info)
	if err != nil {
		log.Fatalf("failed to parse response from Yandex DNS API service %v\n", err)
	}

	return info
}

func getDomainInfo(conf *config) *domainInfo {
	url := fmt.Sprintf(getDomainInfoURLTemplate, conf.Token, conf.Domain)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("request to the Yandex DNS API service failed: %v\n", err)
	}
	defer closeResource(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("request to the Yandex DNS API service failed: %v\n", err)
	}

	info := parseDomainInfoData(body)
	verifyDomainInfoData(info, conf)

	return info
}

func verifyDomainInfoData(info *domainInfo, conf *config) {
	if info.Delegated == nil {
		log.Fatalf("domain is not delegated\n")
	}

	if info.Error != "ok" {
		log.Fatalf("invalid status while calling 'get_domain_records' Yandex DNS API command: %v\n", info.Error)
	}

	if info.Name != conf.Domain {
		log.Fatalf("invariand failed: %s != %s\n", info.Name, conf.Domain)
	}

	if len(info.Records) == 0 {
		log.Fatalf("empty response\n")
	}
}
