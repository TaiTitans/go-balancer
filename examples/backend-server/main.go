package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	port = flag.Int("port", 8081, "Port to listen on")
	name = flag.String("name", "", "Server name (optional)")
)

func main() {
	flag.Parse()

	// If name not provided, use port number
	if *name == "" {
		*name = fmt.Sprintf("Backend-%d", *port)
	}

	mux := http.NewServeMux()

	// Main handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s from %s", *name, r.Method, r.URL.Path, r.RemoteAddr)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Backend-Server", *name)

		response := fmt.Sprintf(`{
  "server": "%s",
  "port": %d,
  "timestamp": "%s",
  "path": "%s",
  "method": "%s",
  "remote_addr": "%s"
}`, *name, *port, time.Now().Format(time.RFC3339), r.URL.Path, r.Method, r.RemoteAddr)

		fmt.Fprint(w, response)
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","server":"%s","port":%d,"timestamp":"%s"}`,
			*name, *port, time.Now().Format(time.RFC3339))
	})

	// Status endpoint
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","server":"%s","port":%d}`, *name, *port)
	})

	addr := fmt.Sprintf(":%d", *port)

	log.Printf("╔════════════════════════════════════════╗")
	log.Printf("║   Backend Server                       ║")
	log.Printf("╚════════════════════════════════════════╝")
	log.Printf("Name:    %s", *name)
	log.Printf("Port:    %d", *port)
	log.Printf("Health:  http://localhost:%d/health", *port)
	log.Printf("════════════════════════════════════════")

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
