package scheduler

// Prevents hammering the same domain.
// Tracks when each domain was last crawled and enforces a minimum delay between requests.

import (
	"sync"
	"time"
)

type PolitenessEnforcer struct {
	mu sync.Mutex
	lastAccess map[string]time.Time
	minDelay time.Duration
}

func NewPolitenessEnforcer(minDelay time.Duration) *PolitenessEnforcer {
	return &PolitenessEnforcer{
		lastAccess: make(map[string]time.Time),
		minDelay: minDelay,
	}
}

func (p *PolitenessEnforcer) CanCrawl(domain string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	lastTime, exists := p.lastAccess[domain]
	if !exists {
		return true
	}
	return time.Since(lastTime) >= p.minDelay
}

func (p *PolitenessEnforcer) RecordAccess(domain string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lastAccess[domain] = time.Now()
}
