package models
import (
	"time"
)
type CrawlTask struct{
	URL string `json:"url"`
	Depth int `json:"depth"`
	JobID string `json:"job_id"`
}

type CrawlResult struct{
	URL string `json:"url"`
	StatusCode int `json:"status_code"`
	Title string `json:"title"`
	Body string `json:"body"`
	Links []string `json:"links"`
	CrawledAt time.Time `json:"crawled_at"`
	ContentLength int `json:"content_length"`
	Depth int `json:"depth"`
	JobID string `json:"job_id"`
}

type FrontierURL struct{
	ID int `json:"id"`
	URL string `json:"url"`
	Domain string `json:"domain"`
	Depth int `json:"depth"`
	Status string `json:"status"`
	Priority int `json:"priority"`
	StatusCode *int `json:"status_code"`
	ErrorMsg *string `json:"error_msg"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CrawlStats struct{
	TotalURLs int `json:"total_urls"`
	Pending int `json:"pending"`
	InProgress int `json:"in_progress"`
	Crawled int `json:"crawled"`
	Failed int `json:"failed"`
}

type PageDocument struct{
	URL string `json:"url"`
	StatusCode int `json:"status_code"`
	Title string `json:"title"`
	Body string `json:"body"`
	Links []string `json:"links"`
	CrawledAt time.Time `json:"crawled_at"`
	ContentLength int `json:"content_length"`
	Depth int `json:"depth"`
	JobID string `json:"job_id"`
	Domain string `json:"domain"`
	IndexedAt time.Time `json:"indexed_at"`
}

type SeedRequest struct{
	URLs []string `json:"urls"`
}

