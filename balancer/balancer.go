package balancer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/TaiTitans/go-balancer/backend"
	"github.com/TaiTitans/go-balancer/healthcheck"
	"github.com/TaiTitans/go-balancer/strategy"
)

// LoadBalancer represents the main load balancer
type LoadBalancer struct {
	backends      []*backend.Backend
	strategy      strategy.Strategy
	healthChecker *healthcheck.HealthChecker
	mu            sync.RWMutex
}

// Config holds the load balancer configuration
type Config struct {
	BackendURLs         []string
	Strategy            strategy.Strategy
	HealthCheckInterval time.Duration
	HealthCheckTimeout  time.Duration
}

// NewLoadBalancer creates a new load balancer instance
func NewLoadBalancer(config Config) (*LoadBalancer, error) {
	if len(config.BackendURLs) == 0 {
		return nil, fmt.Errorf("no backend URLs provided")
	}

	if config.Strategy == nil {
		return nil, fmt.Errorf("no strategy provided")
	}

	// Set default values
	if config.HealthCheckInterval == 0 {
		config.HealthCheckInterval = 10 * time.Second
	}
	if config.HealthCheckTimeout == 0 {
		config.HealthCheckTimeout = 5 * time.Second
	}

	// Create backends
	backends := make([]*backend.Backend, 0, len(config.BackendURLs))
	for _, urlStr := range config.BackendURLs {
		b, err := backend.NewBackend(urlStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create backend for %s: %w", urlStr, err)
		}
		backends = append(backends, b)
	}

	lb := &LoadBalancer{
		backends: backends,
		strategy: config.Strategy,
	}

	// Create health checker
	lb.healthChecker = healthcheck.NewHealthChecker(
		backends,
		config.HealthCheckInterval,
		config.HealthCheckTimeout,
	)

	return lb, nil
}

// Start starts the load balancer and health checker
func (lb *LoadBalancer) Start(ctx context.Context) {
	log.Printf("Starting load balancer with strategy: %s", lb.strategy.Name())
	log.Printf("Managing %d backend(s)", len(lb.backends))

	go lb.healthChecker.Start(ctx)
}

// ServeHTTP implements the http.Handler interface
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Select a backend
	selectedBackend := lb.strategy.SelectBackend(lb.backends)
	if selectedBackend == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		log.Println("No available backends")
		return
	}

	// Increment connection count
	selectedBackend.IncrementConnections()
	defer selectedBackend.DecrementConnections()

	log.Printf("Forwarding request to %s (active connections: %d)",
		selectedBackend.GetURL(), selectedBackend.GetConnections())

	// Create reverse proxy
	proxy := lb.createProxy(selectedBackend)
	proxy.ServeHTTP(w, r)
}

// createProxy creates a reverse proxy for the given backend
func (lb *LoadBalancer) createProxy(b *backend.Backend) *httputil.ReverseProxy {
	targetURL := b.GetURL()

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Custom director to modify the request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
		req.URL.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
	}

	// Error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error for backend %s: %v", targetURL, err)
		b.SetAlive(false)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	return proxy
}

// GetBackends returns all backends
func (lb *LoadBalancer) GetBackends() []*backend.Backend {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return lb.backends
}

// GetStrategy returns the current strategy
func (lb *LoadBalancer) GetStrategy() strategy.Strategy {
	return lb.strategy
}

// SetStrategy sets a new load balancing strategy
func (lb *LoadBalancer) SetStrategy(s strategy.Strategy) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.strategy = s
	log.Printf("Strategy changed to: %s", s.Name())
}

// GetStats returns statistics about the backends
func (lb *LoadBalancer) GetStats() map[string]interface{} {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	stats := make(map[string]interface{})
	backendStats := make([]map[string]interface{}, 0, len(lb.backends))

	totalAlive := 0
	for _, b := range lb.backends {
		alive := b.IsAlive()
		if alive {
			totalAlive++
		}

		backendStats = append(backendStats, map[string]interface{}{
			"url":          b.GetURL().String(),
			"alive":        alive,
			"connections":  b.GetConnections(),
			"responseTime": b.GetResponseTime().String(),
		})
	}

	stats["strategy"] = lb.strategy.Name()
	stats["totalBackends"] = len(lb.backends)
	stats["aliveBackends"] = totalAlive
	stats["backends"] = backendStats

	return stats
}

// HandleStats returns an HTTP handler for stats endpoint
func (lb *LoadBalancer) HandleStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := lb.GetStats()

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "Load Balancer Statistics\n")
		fmt.Fprintf(w, "========================\n\n")
		fmt.Fprintf(w, "Strategy: %s\n", stats["strategy"])
		fmt.Fprintf(w, "Total Backends: %d\n", stats["totalBackends"])
		fmt.Fprintf(w, "Alive Backends: %d\n\n", stats["aliveBackends"])

		fmt.Fprintf(w, "Backend Details:\n")
		fmt.Fprintf(w, "----------------\n")

		if backends, ok := stats["backends"].([]map[string]interface{}); ok {
			for i, b := range backends {
				fmt.Fprintf(w, "%d. %s\n", i+1, b["url"])
				fmt.Fprintf(w, "   Status: ")
				if b["alive"].(bool) {
					fmt.Fprintf(w, "✓ Alive\n")
				} else {
					fmt.Fprintf(w, "✗ Down\n")
				}
				fmt.Fprintf(w, "   Active Connections: %d\n", b["connections"])
				fmt.Fprintf(w, "   Response Time: %s\n\n", b["responseTime"])
			}
		}
	}
}
