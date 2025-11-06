package strategy

import (
	"testing"

	"github.com/TaiTitans/go-balancer/backend"
)

func createTestBackends(count int) []*backend.Backend {
	backends := make([]*backend.Backend, count)
	for i := 0; i < count; i++ {
		b, _ := backend.NewBackend("http://localhost:808" + string(rune('0'+i)))
		backends[i] = b
	}
	return backends
}

func TestRoundRobin(t *testing.T) {
	strategy := NewRoundRobin()
	backends := createTestBackends(3)

	if strategy.Name() != "RoundRobin" {
		t.Errorf("Expected strategy name 'RoundRobin', got '%s'", strategy.Name())
	}

	// Test round-robin selection
	selected := make(map[*backend.Backend]int)
	for i := 0; i < 9; i++ {
		b := strategy.SelectBackend(backends)
		if b == nil {
			t.Error("SelectBackend returned nil")
		}
		selected[b]++
	}

	// Each backend should be selected 3 times
	for _, count := range selected {
		if count != 3 {
			t.Errorf("Expected each backend to be selected 3 times, got %d", count)
		}
	}
}

func TestRoundRobin_EmptyBackends(t *testing.T) {
	strategy := NewRoundRobin()
	b := strategy.SelectBackend([]*backend.Backend{})
	if b != nil {
		t.Error("SelectBackend should return nil for empty backends")
	}
}

func TestRoundRobin_DeadBackends(t *testing.T) {
	strategy := NewRoundRobin()
	backends := createTestBackends(3)

	// Mark all backends as dead
	for _, b := range backends {
		b.SetAlive(false)
	}

	selected := strategy.SelectBackend(backends)
	if selected != nil {
		t.Error("SelectBackend should return nil when all backends are dead")
	}
}

func TestLeastConnections(t *testing.T) {
	strategy := NewLeastConnections()
	backends := createTestBackends(3)

	if strategy.Name() != "LeastConnections" {
		t.Errorf("Expected strategy name 'LeastConnections', got '%s'", strategy.Name())
	}

	// Set different connection counts
	backends[0].IncrementConnections()
	backends[0].IncrementConnections()
	backends[1].IncrementConnections()

	selected := strategy.SelectBackend(backends)
	if selected != backends[2] {
		t.Error("LeastConnections should select backend with 0 connections")
	}
}

func TestLeastConnections_EmptyBackends(t *testing.T) {
	strategy := NewLeastConnections()
	b := strategy.SelectBackend([]*backend.Backend{})
	if b != nil {
		t.Error("SelectBackend should return nil for empty backends")
	}
}

func TestRandom(t *testing.T) {
	strategy := NewRandom()
	backends := createTestBackends(3)

	if strategy.Name() != "Random" {
		t.Errorf("Expected strategy name 'Random', got '%s'", strategy.Name())
	}

	// Test that random selection works
	selected := make(map[*backend.Backend]bool)
	for i := 0; i < 100; i++ {
		b := strategy.SelectBackend(backends)
		if b == nil {
			t.Error("SelectBackend returned nil")
		}
		selected[b] = true
	}

	// With 100 iterations, we should have selected all backends at least once
	if len(selected) != len(backends) {
		t.Errorf("Expected all backends to be selected, got %d/%d", len(selected), len(backends))
	}
}

func TestRandom_EmptyBackends(t *testing.T) {
	strategy := NewRandom()
	b := strategy.SelectBackend([]*backend.Backend{})
	if b != nil {
		t.Error("SelectBackend should return nil for empty backends")
	}
}
