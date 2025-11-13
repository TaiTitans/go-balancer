package balancer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/TaiTitans/go-balancer/backend"
	constants "github.com/TaiTitans/go-balancer/const"
	"github.com/TaiTitans/go-balancer/healthcheck"
	"github.com/TaiTitans/go-balancer/strategy"
)

// LoadBalancer represents the main load balancer
type LoadBalancer struct {
	backends      []*backend.Backend
	strategy      strategy.Strategy
	healthChecker *healthcheck.HealthChecker
	mu            sync.RWMutex
	metrics       *Metrics
}

// Metrics tracks load balancer performance
type Metrics struct {
	TotalRequests  int64
	FailedRequests int64
	TotalBytes     int64
	mu             sync.RWMutex
	StartTime      time.Time
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
		config.HealthCheckInterval = constants.DefaultHealthCheckInterval
	}
	if config.HealthCheckTimeout == 0 {
		config.HealthCheckTimeout = constants.DefaultHealthCheckTimeout
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
		metrics: &Metrics{
			StartTime: time.Now(),
		},
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
	atomic.AddInt64(&lb.metrics.TotalRequests, 1)

	// Select a backend using the strategy
	lb.mu.RLock()
	selectedBackend := lb.strategy.SelectBackend(lb.backends)
	lb.mu.RUnlock()

	if selectedBackend == nil {
		atomic.AddInt64(&lb.metrics.FailedRequests, 1)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		log.Println("No available backends")
		return
	}

	log.Printf("Forwarding request to %s (active connections: %d)",
		selectedBackend.GetURL(), selectedBackend.GetConnections())

	// Use the backend's Serve method which already has ReverseProxy configured
	selectedBackend.Serve(w, r)
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
	totalConnections := 0
	for _, b := range lb.backends {
		alive := b.IsAlive()
		if alive {
			totalAlive++
		}
		connections := b.GetConnections()
		totalConnections += connections

		backendStats = append(backendStats, map[string]interface{}{
			"url":          b.GetURL().String(),
			"alive":        alive,
			"connections":  connections,
			"responseTime": b.GetResponseTime().String(),
			"failCount":    b.GetFailCount(),
		})
	}

	uptime := time.Since(lb.metrics.StartTime)
	totalReqs := atomic.LoadInt64(&lb.metrics.TotalRequests)
	failedReqs := atomic.LoadInt64(&lb.metrics.FailedRequests)

	stats["strategy"] = lb.strategy.Name()
	stats["totalBackends"] = len(lb.backends)
	stats["aliveBackends"] = totalAlive
	stats["totalConnections"] = totalConnections
	stats["totalRequests"] = totalReqs
	stats["failedRequests"] = failedReqs
	stats["successRate"] = calculateSuccessRate(totalReqs, failedReqs)
	stats["uptime"] = uptime.String()
	stats["backends"] = backendStats

	return stats
}

func calculateSuccessRate(total, failed int64) string {
	if total == 0 {
		return "N/A"
	}
	rate := float64(total-failed) / float64(total) * 100
	return fmt.Sprintf("%.2f%%", rate)
}

// HandleStats returns an HTTP handler for stats endpoint
func (lb *LoadBalancer) HandleStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := lb.GetStats()

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "╔════════════════════════════════════════╗\n")
		fmt.Fprintf(w, "║   Load Balancer Statistics             ║\n")
		fmt.Fprintf(w, "╚════════════════════════════════════════╝\n\n")

		fmt.Fprintf(w, "Strategy:         %s\n", stats["strategy"])
		fmt.Fprintf(w, "Uptime:           %s\n", stats["uptime"])
		fmt.Fprintf(w, "Total Backends:   %d\n", stats["totalBackends"])
		fmt.Fprintf(w, "Alive Backends:   %d\n", stats["aliveBackends"])
		fmt.Fprintf(w, "Total Requests:   %d\n", stats["totalRequests"])
		fmt.Fprintf(w, "Failed Requests:  %d\n", stats["failedRequests"])
		fmt.Fprintf(w, "Success Rate:     %s\n", stats["successRate"])
		fmt.Fprintf(w, "Active Connections: %d\n\n", stats["totalConnections"])

		fmt.Fprintf(w, "Backend Details:\n")
		fmt.Fprintf(w, "════════════════════════════════════════\n")

		if backends, ok := stats["backends"].([]map[string]interface{}); ok {
			for i, b := range backends {
				fmt.Fprintf(w, "\n[%d] %s\n", i+1, b["url"])
				if b["alive"].(bool) {
					fmt.Fprintf(w, "    Status:       ✓ Healthy\n")
				} else {
					fmt.Fprintf(w, "    Status:       ✗ Down\n")
				}
				fmt.Fprintf(w, "    Connections:  %d\n", b["connections"])
				fmt.Fprintf(w, "    Response Time: %s\n", b["responseTime"])
				fmt.Fprintf(w, "    Fail Count:   %d\n", b["failCount"])
			}
		}

		fmt.Fprintf(w, "\n════════════════════════════════════════\n")
	}
}
