package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/TaiTitans/go-balancer/balancer"
	"github.com/TaiTitans/go-balancer/middleware"
	"github.com/TaiTitans/go-balancer/strategy"
)

var (
	port           = flag.Int("port", 8080, "Load balancer port")
	backendsFlag   = flag.String("backends", "http://localhost:8081,http://localhost:8082,http://localhost:8083", "Comma-separated list of backend URLs")
	strategyFlag   = flag.String("strategy", "roundrobin", "Load balancing strategy (roundrobin, leastconnections, random)")
	healthInterval = flag.Duration("health-interval", 10*time.Second, "Health check interval")
	healthTimeout  = flag.Duration("health-timeout", 5*time.Second, "Health check timeout")
)

func main() {
	flag.Parse()

	// Parse backend URLs
	backendURLs := parseBackendURLs(*backendsFlag)
	if len(backendURLs) == 0 {
		log.Fatal("No backend URLs provided")
	}

	// Select strategy
	var strat strategy.Strategy
	switch strings.ToLower(*strategyFlag) {
	case "roundrobin":
		strat = strategy.NewRoundRobin()
	case "leastconnections":
		strat = strategy.NewLeastConnections()
	case "random":
		strat = strategy.NewRandom()
	default:
		log.Fatalf("Unknown strategy: %s", *strategyFlag)
	}

	// Configure the load balancer
	config := balancer.Config{
		BackendURLs:         backendURLs,
		Strategy:            strat,
		HealthCheckInterval: *healthInterval,
		HealthCheckTimeout:  *healthTimeout,
	}

	// Create load balancer
	lb, err := balancer.NewLoadBalancer(config)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the load balancer
	lb.Start(ctx)

	// Create HTTP server with middleware
	mux := http.NewServeMux()
	mux.Handle("/", lb)
	mux.Handle("/stats", lb.HandleStats())
	mux.HandleFunc("/health", healthHandler)

	// Apply middleware
	handler := middleware.Chain(
		mux,
		middleware.Logger,
		middleware.Recovery,
		middleware.CORS,
	)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("╔════════════════════════════════════════╗")
		log.Printf("║   Go Load Balancer                     ║")
		log.Printf("╚════════════════════════════════════════╝")
		log.Printf("Port:          %d", *port)
		log.Printf("Strategy:      %s", strat.Name())
		log.Printf("Backends:      %d", len(backendURLs))
		log.Printf("Health Check:  %v", *healthInterval)
		log.Printf("")
		log.Printf("Endpoints:")
		log.Printf("  - Load Balancer: http://localhost:%d/", *port)
		log.Printf("  - Statistics:    http://localhost:%d/stats", *port)
		log.Printf("  - Health:        http://localhost:%d/health", *port)
		log.Printf("")
		log.Printf("Backends:")
		for i, url := range backendURLs {
			log.Printf("  [%d] %s", i+1, url)
		}
		log.Printf("════════════════════════════════════════")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\nShutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func parseBackendURLs(backends string) []string {
	if backends == "" {
		return nil
	}

	urls := strings.Split(backends, ",")
	result := make([]string, 0, len(urls))

	for _, url := range urls {
		trimmed := strings.TrimSpace(url)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}
