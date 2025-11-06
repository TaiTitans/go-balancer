package healthcheck

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/TaiTitans/go-balancer/backend"
)

// HealthChecker performs health checks on backends
type HealthChecker struct {
	backends []*backend.Backend
	interval time.Duration
	timeout  time.Duration
	client   *http.Client
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(backends []*backend.Backend, interval, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		backends: backends,
		interval: interval,
		timeout:  timeout,
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: timeout,
				}).DialContext,
			},
		},
	}
}

// Start begins the health check loop
func (hc *HealthChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	// Perform initial health check
	hc.checkAll()

	for {
		select {
		case <-ctx.Done():
			log.Println("Health checker stopped")
			return
		case <-ticker.C:
			hc.checkAll()
		}
	}
}

// checkAll checks all backends
func (hc *HealthChecker) checkAll() {
	for _, b := range hc.backends {
		go hc.check(b)
	}
}

// check performs a health check on a single backend
func (hc *HealthChecker) check(b *backend.Backend) {
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, b.GetURL().String(), nil)
	if err != nil {
		b.SetAlive(false)
		log.Printf("Failed to create request for %s: %v", b.GetURL(), err)
		return
	}

	resp, err := hc.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		b.SetAlive(false)
		log.Printf("Backend %s is down: %v", b.GetURL(), err)
		return
	}
	defer resp.Body.Close()

	// Consider 2xx and 3xx as healthy
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		b.SetAlive(true)
		b.UpdateResponseTime(duration)
		log.Printf("Backend %s is healthy (response time: %v)", b.GetURL(), duration)
	} else {
		b.SetAlive(false)
		log.Printf("Backend %s returned status %d", b.GetURL(), resp.StatusCode)
	}
}
