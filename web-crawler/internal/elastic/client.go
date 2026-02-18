package elastic

import (
	"os"
	"log"
	"strings"
	"github.com/elastic/go-elasticsearch/v9"
)

func InitClient()*elasticsearch.Client{
	ELASTICSEARCH_URL := os.Getenv("ELASTICSEARCH_URL")
	cfg := elasticsearch.Config{
		Addresses: []string{ELASTICSEARCH_URL},
	}
	es,err := elasticsearch.NewClient(cfg)
	if err != nil{
		log.Fatal("failed to setup elasticsearch")
	}
	_,err1:=es.Info()
	if err1 != nil{
		log.Fatal("failed to setup elasticsearch")
	}
	mappings := `
	{
		"mappings": {
			"properties": {
				"url": {"type":"keyword"},
				"domain": {"type":"keyword"},
				"title": {"type":"text"},
				"body": {"type":"text"},
				"links": {"type":"keyword"},
				"status_code": {"type":"integer"},
				"content_length": {"type":"integer"},
				"depth": {"type":"integer"},
				"crawled_at": {"type":"date"},
				"indexed_at": {"type":"date"}
			}
		}
	}
	`
	es.Indices.Create("crawled_pages",es.Indices.Create.WithBody(strings.NewReader(mappings)))
	return es
}

