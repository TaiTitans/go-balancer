package strategy

import (
	"math/rand"
	"time"

	"github.com/TaiTitans/go-balancer/backend"
)

// Random implements random load balancing strategy
type Random struct {
	rng *rand.Rand
}

// NewRandom creates a new random strategy
func NewRandom() *Random {
	return &Random{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SelectBackend selects a random backend from alive backends
func (r *Random) SelectBackend(backends []*backend.Backend) *backend.Backend {
	if len(backends) == 0 {
		return nil
	}

	// Find alive backends
	var aliveBackends []*backend.Backend
	for _, b := range backends {
		if b.IsAlive() {
			aliveBackends = append(aliveBackends, b)
		}
	}

	if len(aliveBackends) == 0 {
		return nil
	}

	// Select random backend
	idx := r.rng.Intn(len(aliveBackends))
	return aliveBackends[idx]
}

// Name returns the strategy name
func (r *Random) Name() string {
	return "Random"
}
