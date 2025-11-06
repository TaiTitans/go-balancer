package strategy

import (
	"github.com/TaiTitans/go-balancer/backend"
)

// Strategy defines the interface for load balancing strategies
type Strategy interface {
	// SelectBackend selects a backend from the pool
	SelectBackend(backends []*backend.Backend) *backend.Backend
	// Name returns the name of the strategy
	Name() string
}
