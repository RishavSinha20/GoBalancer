package balancer

import (
	"net/url"
	"sync"
)

type Backend struct {
	URL         *url.URL
	alive       bool
	connections int
	mu          sync.RWMutex
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
	defer b.mu.Unlock()
	return b.alive
}

func (b *Backend) IncrementConnections() {
	b.mu.Lock()
	b.connections++
	b.mu.Unlock()
}

func (b *Backend) DecrementConnections() {
	b.mu.Lock()
	b.connections--
	b.mu.Unlock()
}
func (b *Backend) GetConnections() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.connections
}
