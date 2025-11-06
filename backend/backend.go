package backend

import (
	"net/url"
	"sync"
	"time"
)

// Backend represents a backend server
type Backend struct {
	URL          *url.URL
	Alive        bool
	mu           sync.RWMutex
	Connections  int
	ResponseTime time.Duration
}

// NewBackend creates a new backend instance
func NewBackend(urlStr string) (*Backend, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &Backend{
		URL:   u,
		Alive: true,
	}, nil
}

// SetAlive sets the alive status of the backend
func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
}

// IsAlive returns the alive status of the backend
func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}

// GetURL returns the backend URL
func (b *Backend) GetURL() *url.URL {
	return b.URL
}

// IncrementConnections increments the connection count
func (b *Backend) IncrementConnections() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Connections++
}

// DecrementConnections decrements the connection count
func (b *Backend) DecrementConnections() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.Connections > 0 {
		b.Connections--
	}
}

// GetConnections returns the current connection count
func (b *Backend) GetConnections() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Connections
}

// UpdateResponseTime updates the response time of the backend
func (b *Backend) UpdateResponseTime(duration time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ResponseTime = duration
}

// GetResponseTime returns the response time of the backend
func (b *Backend) GetResponseTime() time.Duration {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.ResponseTime
}
