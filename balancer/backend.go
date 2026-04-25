package balancer

import (
	"net/url"
	"sync"
)

type Backend struct {
	URL   *url.URL
	alive bool
	mu    sync.RWMutex
}

func NewBackend(rawURL string) *Backend {
	parsed, _ := url.Parse(rawURL)

	return &Backend{
		URL:   parsed,
		alive: true,
	}
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	b.alive = alive
	b.mu.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.mu.RLock()

	defer b.mu.RUnlock()
	return b.alive
}
