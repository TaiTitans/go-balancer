package backend

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Backend represents a backend server
type Backend struct {
	URL          *url.URL
	Alive        bool
	mu           sync.RWMutex
	Connections  int32
	ResponseTime time.Duration
	ReverseProxy *httputil.ReverseProxy
	FailCount    int32
	LastCheck    time.Time
}

// Serve handles the HTTP request by forwarding it to the backend server
func (b *Backend) Serve(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	atomic.AddInt32(&b.Connections, 1)
	defer func() {
		atomic.AddInt32(&b.Connections, -1)
		b.UpdateResponseTime(time.Since(start))
	}()
	b.ReverseProxy.ServeHTTP(w, r)
}

// ServerPool manages a pool of backend servers
type ServerPool struct {
	backends []*Backend
	current  uint64
	mu       sync.RWMutex
}

// NewBackend creates a new backend instance with enhanced configuration
func NewBackend(urlStr string) (*Backend, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	b := &Backend{
		URL:       u,
		Alive:     true,
		LastCheck: time.Now(),
	}

	// Create reverse proxy with custom configuration
	rp := httputil.NewSingleHostReverseProxy(u)

	// Custom director to properly forward requests
	originalDirector := rp.Director
	rp.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = u.Host
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Origin-Host", u.Host)
		req.Header.Set("X-Forwarded-Proto", req.URL.Scheme)
	}

	// Error handler with automatic retry and failure tracking
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[Backend Error] %s: %v", u, err)
		atomic.AddInt32(&b.FailCount, 1)
		b.SetAlive(false)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	// Custom response modifier for logging
	rp.ModifyResponse = func(resp *http.Response) error {
		// Reset fail count on successful response
		if resp.StatusCode < 500 {
			atomic.StoreInt32(&b.FailCount, 0)
		}
		return nil
	}

	b.ReverseProxy = rp

	return b, nil
}

// SetAlive sets the alive status of the backend
func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
}

// IsAlive returns the alive status of the backend
func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}

// GetURL returns the backend URL
func (b *Backend) GetURL() *url.URL {
	return b.URL
}

// IncrementConnections increments the connection count atomically
func (b *Backend) IncrementConnections() {
	atomic.AddInt32(&b.Connections, 1)
}

// DecrementConnections decrements the connection count atomically
func (b *Backend) DecrementConnections() {
	atomic.AddInt32(&b.Connections, -1)
}

// GetConnections returns the current connection count
func (b *Backend) GetConnections() int {
	return int(atomic.LoadInt32(&b.Connections))
}

// GetFailCount returns the current failure count
func (b *Backend) GetFailCount() int {
	return int(atomic.LoadInt32(&b.FailCount))
}

// ResetFailCount resets the failure count
func (b *Backend) ResetFailCount() {
	atomic.StoreInt32(&b.FailCount, 0)
}

// SetLastCheck updates the last health check time
func (b *Backend) SetLastCheck(t time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.LastCheck = t
}

// GetLastCheck returns the last health check time
func (b *Backend) GetLastCheck() time.Time {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.LastCheck
}

// UpdateResponseTime updates the response time of the backend
func (b *Backend) UpdateResponseTime(duration time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ResponseTime = duration
}

// GetResponseTime returns the response time of the backend
func (b *Backend) GetResponseTime() time.Duration {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.ResponseTime
}

// NewServerPool creates a new server pool
func NewServerPool() *ServerPool {
	return &ServerPool{
		backends: make([]*Backend, 0),
	}
}

// AddBackend adds a backend to the pool
func (sp *ServerPool) AddBackend(b *Backend) {
	sp.backends = append(sp.backends, b)
}

// GetBackends returns all backends in the pool
func (sp *ServerPool) GetBackends() []*Backend {
	return sp.backends
}

// NextIndex atomically increases the current index and returns the next index
func (sp *ServerPool) NextIndex() int {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	if len(sp.backends) == 0 {
		return 0
	}
	return int(atomic.AddUint64(&sp.current, 1) % uint64(len(sp.backends)))
}

// GetNextPeer returns the next available backend to handle the request
func (sp *ServerPool) GetNextPeer() *Backend {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	if len(sp.backends) == 0 {
		return nil
	}

	next := sp.NextIndex()
	l := len(sp.backends)

	for i := 0; i < l; i++ {
		idx := (next + i) % l
		if sp.backends[idx].IsAlive() {
			return sp.backends[idx]
		}
	}

	return nil
}

// GetAliveBackends returns a list of all alive backends
func (sp *ServerPool) GetAliveBackends() []*Backend {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	alive := make([]*Backend, 0)
	for _, b := range sp.backends {
		if b.IsAlive() {
			alive = append(alive, b)
		}
	}
	return alive
}

// GetBackendCount returns the total number of backends
func (sp *ServerPool) GetBackendCount() int {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	return len(sp.backends)
}

// MarkBackendStatus changes the status of a backend identified by URL
func (sp *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range sp.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}
