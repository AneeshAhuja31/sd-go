package scheduler

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"web-crawler/internal/bloom"
	"web-crawler/internal/db"
	"web-crawler/internal/hashring"
	"web-crawler/internal/models"
	"web-crawler/internal/queue"
)

type Scheduler struct {
	DB *sql.DB
	Ch *amqp.Channel
	Bloom *bloom.BloomFilter
	Ring *hashring.HashRing
	Politeness *PolitenessEnforcer
}

func NewScheduler(pg *sql.DB, ch *amqp.Channel, bf *bloom.BloomFilter, ring *hashring.HashRing, politeness *PolitenessEnforcer) *Scheduler {
	return &Scheduler{
		DB: pg,
		Ch: ch,
		Bloom: bf,
		Ring: ring,
		Politeness: politeness,
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(2 * time.Second)
	log.Println("Scheduler started, ticking every 2s")

	for range ticker.C {
		urls, err := db.DequeueURLs(s.DB, 50)
		if err != nil {
			log.Println("Scheduler dequeue error: ", err)
			continue
		}
		if len(urls) == 0 {
			continue
		}

		for _, frontierURL := range urls {
			if s.Bloom.MightContain(frontierURL.URL) {
				continue
			}

			parsed, err := url.Parse(frontierURL.URL)
			if err != nil {
				log.Println("Error parsing URL: ", frontierURL.URL)
				continue
			}
			domain := parsed.Hostname()

			if !s.Politeness.CanCrawl(domain) {
				continue
			}

			workerID := s.Ring.GetWorker(domain)

			task := models.CrawlTask{
				URL: frontierURL.URL,
				Depth: frontierURL.Depth,
				JobID: fmt.Sprintf("job-%d", frontierURL.ID),
			}
			err = queue.PublishCrawlTask(s.Ch, workerID, task)
			if err != nil {
				log.Println("Error publishing task: ", err)
				continue
			}

			s.Bloom.Add(frontierURL.URL)
			s.Politeness.RecordAccess(domain)

			log.Printf("Dispatched %s → worker %d", frontierURL.URL, workerID)
		}
	}
}
