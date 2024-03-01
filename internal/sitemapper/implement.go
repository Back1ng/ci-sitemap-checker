package sitemapper

import (
	"crypto/tls"
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

type sitemapper struct {
}

func New() SitemapParser {
	return &sitemapper{}
}

func (s *sitemapper) Get(url string) Sitemap {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var parsedXml Sitemap
	xml.Unmarshal(body, &parsedXml)

	return parsedXml
}
