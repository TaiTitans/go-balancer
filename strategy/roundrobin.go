package strategy

import (
	"sync/atomic"

	"github.com/TaiTitans/go-balancer/backend"
)

// RoundRobin implements round-robin load balancing strategy
type RoundRobin struct {
	current uint64
}

// NewRoundRobin creates a new round-robin strategy
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{current: 0}
}

// SelectBackend selects the next backend in round-robin fashion
func (rr *RoundRobin) SelectBackend(backends []*backend.Backend) *backend.Backend {
	if len(backends) == 0 {
		return nil
	}

	// Find alive backends
	aliveBackends := []*backend.Backend{}
	for _, b := range backends {
		if b.IsAlive() {
			aliveBackends = append(aliveBackends, b)
		}
	}

	if len(aliveBackends) == 0 {
		return nil
	}

	// Get next backend using atomic operation
	next := atomic.AddUint64(&rr.current, 1)
	return aliveBackends[(int(next)-1)%len(aliveBackends)]
}

// Name returns the strategy name
func (rr *RoundRobin) Name() string {
	return RoundRobinStrategy
}
