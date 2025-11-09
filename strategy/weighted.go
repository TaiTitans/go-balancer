package strategy

import (
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/TaiTitans/go-balancer/backend"
)

// WeightedRoundRobin implements weighted round-robin load balancing strategy
type WeightedRoundRobin struct {
	current uint64
	weights map[*backend.Backend]int
	rng     *rand.Rand
}

// NewWeightedRoundRobin creates a new weighted round-robin strategy
func NewWeightedRoundRobin(weights map[*backend.Backend]int) *WeightedRoundRobin {
	return &WeightedRoundRobin{
		current: 0,
		weights: weights,
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SelectBackend selects a backend based on weighted round-robin
func (wrr *WeightedRoundRobin) SelectBackend(backends []*backend.Backend) *backend.Backend {
	if len(backends) == 0 {
		return nil
	}

	// Find alive backends with weights
	weightedBackends := []*backend.Backend{}
	totalWeight := 0

	for _, b := range backends {
		if b.IsAlive() {
			weight := 1
			if w, ok := wrr.weights[b]; ok {
				weight = w
			}

			// Add backend multiple times based on weight
			for i := 0; i < weight; i++ {
				weightedBackends = append(weightedBackends, b)
				totalWeight++
			}
		}
	}

	if len(weightedBackends) == 0 {
		return nil
	}

	// Select using round-robin on weighted list
	next := atomic.AddUint64(&wrr.current, 1)
	return weightedBackends[(int(next)-1)%len(weightedBackends)]
}

// Name returns the strategy name
func (wrr *WeightedRoundRobin) Name() string {
	return WeightedRoundRobinStrategy
}

// IPHash implements IP hash load balancing strategy
type IPHash struct{}

// NewIPHash creates a new IP hash strategy
func NewIPHash() *IPHash {
	return &IPHash{}
}

// SelectBackend selects a backend based on client IP hash
func (ih *IPHash) SelectBackend(backends []*backend.Backend) *backend.Backend {
	// Note: This requires access to request, so it's a simplified version
	// In practice, you'd need to pass the request to calculate IP hash
	if len(backends) == 0 {
		return nil
	}

	aliveBackends := []*backend.Backend{}
	for _, b := range backends {
		if b.IsAlive() {
			aliveBackends = append(aliveBackends, b)
		}
	}

	if len(aliveBackends) == 0 {
		return nil
	}

	// Simplified: just return first alive backend
	// Real implementation would hash client IP
	return aliveBackends[0]
}

// Name returns the strategy name
func (ih *IPHash) Name() string {
	return IPHashStrategy
}
