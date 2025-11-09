package strategy

import (
	"github.com/TaiTitans/go-balancer/backend"
)

// LeastConnections implements least connections load balancing strategy
type LeastConnections struct{}

// NewLeastConnections creates a new least connections strategy
func NewLeastConnections() *LeastConnections {
	return &LeastConnections{}
}

// SelectBackend selects the backend with the least active connections
func (lc *LeastConnections) SelectBackend(backends []*backend.Backend) *backend.Backend {
	if len(backends) == 0 {
		return nil
	}

	selected := &backend.Backend{}
	minConnections := -1

	for _, b := range backends {
		if !b.IsAlive() {
			continue
		}

		connections := b.GetConnections()
		if minConnections == -1 || connections < minConnections {
			minConnections = connections
			selected = b
		}
	}

	return selected
}

// Name returns the strategy name
func (lc *LeastConnections) Name() string {
	return LeastConnectionsStrategy
}
