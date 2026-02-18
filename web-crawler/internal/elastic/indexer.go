package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"web-crawler/internal/models"

	"github.com/elastic/go-elasticsearch/v9"
)

func IndexPage(es *elasticsearch.Client, doc models.PageDocument) error{
	doc.IndexedAt = time.Now()
	data,err := json.Marshal(doc)
	if err != nil{
		log.Println("Error converting page document to json: ",err)
		return err
	}
	res,err:= es.Index("crawled_pages",bytes.NewReader(data))
	
	if err != nil{
		log.Println("Error indexing document in elasticserch: ",err)
	}
	if res.IsError(){
		return fmt.Errorf("elasticsearch index error: %s", res.String())
	}
	defer res.Body.Close()
	return nil
}