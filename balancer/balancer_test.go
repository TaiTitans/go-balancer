package balancer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TaiTitans/go-balancer/strategy"
)

func TestNewLoadBalancer(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				BackendURLs:         []string{"http://localhost:8081"},
				Strategy:            strategy.NewRoundRobin(),
				HealthCheckInterval: 10 * time.Second,
				HealthCheckTimeout:  5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "no backends",
			config: Config{
				BackendURLs: []string{},
				Strategy:    strategy.NewRoundRobin(),
			},
			wantErr: true,
		},
		{
			name: "no strategy",
			config: Config{
				BackendURLs: []string{"http://localhost:8081"},
			},
			wantErr: true,
		},
		{
			name: "invalid backend URL",
			config: Config{
				BackendURLs: []string{"://invalid"},
				Strategy:    strategy.NewRoundRobin(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb, err := NewLoadBalancer(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLoadBalancer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && lb == nil {
				t.Error("NewLoadBalancer() returned nil")
			}
		})
	}
}

func TestLoadBalancer_ServeHTTP(t *testing.T) {
	// Create test backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend response"))
	}))
	defer backend.Close()

	// Create load balancer
	config := Config{
		BackendURLs:         []string{backend.URL},
		Strategy:            strategy.NewRoundRobin(),
		HealthCheckInterval: 10 * time.Second,
		HealthCheckTimeout:  5 * time.Second,
	}

	lb, err := NewLoadBalancer(config)
	if err != nil {
		t.Fatalf("Failed to create load balancer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lb.Start(ctx)

	// Wait for health check
	time.Sleep(100 * time.Millisecond)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	// Serve request
	lb.ServeHTTP(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "backend response" {
		t.Errorf("Expected body 'backend response', got '%s'", rr.Body.String())
	}
}

func TestLoadBalancer_GetStats(t *testing.T) {
	config := Config{
		BackendURLs:         []string{"http://localhost:8081", "http://localhost:8082"},
		Strategy:            strategy.NewRoundRobin(),
		HealthCheckInterval: 10 * time.Second,
		HealthCheckTimeout:  5 * time.Second,
	}

	lb, err := NewLoadBalancer(config)
	if err != nil {
		t.Fatalf("Failed to create load balancer: %v", err)
	}

	stats := lb.GetStats()

	if stats["totalBackends"] != 2 {
		t.Errorf("Expected 2 backends, got %v", stats["totalBackends"])
	}

	if stats["strategy"] != "RoundRobin" {
		t.Errorf("Expected RoundRobin strategy, got %v", stats["strategy"])
	}
}

func TestLoadBalancer_SetStrategy(t *testing.T) {
	config := Config{
		BackendURLs:         []string{"http://localhost:8081"},
		Strategy:            strategy.NewRoundRobin(),
		HealthCheckInterval: 10 * time.Second,
		HealthCheckTimeout:  5 * time.Second,
	}

	lb, err := NewLoadBalancer(config)
	if err != nil {
		t.Fatalf("Failed to create load balancer: %v", err)
	}

	// Change strategy
	newStrategy := strategy.NewLeastConnections()
	lb.SetStrategy(newStrategy)

	if lb.GetStrategy().Name() != "LeastConnections" {
		t.Errorf("Expected LeastConnections strategy, got %s", lb.GetStrategy().Name())
	}
}

func BenchmarkLoadBalancer_ServeHTTP(b *testing.B) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	config := Config{
		BackendURLs:         []string{backend.URL},
		Strategy:            strategy.NewRoundRobin(),
		HealthCheckInterval: 10 * time.Second,
		HealthCheckTimeout:  5 * time.Second,
	}

	lb, _ := NewLoadBalancer(config)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lb.Start(ctx)

	time.Sleep(100 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		lb.ServeHTTP(rr, req)
	}
}
