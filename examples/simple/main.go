package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TaiTitans/go-balancer/balancer"
	"github.com/TaiTitans/go-balancer/strategy"
)

func main() {
	// Configure the load balancer
	config := balancer.Config{
		BackendURLs: []string{
			"http://localhost:8081",
			"http://localhost:8082",
			"http://localhost:8083",
		},
		Strategy:            strategy.NewRoundRobin(), // Try: NewLeastConnections(), NewRandom()
		HealthCheckInterval: 10 * time.Second,
		HealthCheckTimeout:  5 * time.Second,
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

	// Create HTTP server
	mux := http.NewServeMux()
	mux.Handle("/", lb)
	mux.Handle("/stats", lb.HandleStats())

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Load balancer listening on %s", server.Addr)
		log.Printf("Stats available at http://localhost%s/stats", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
