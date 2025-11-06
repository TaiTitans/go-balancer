package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	port := flag.Int("port", 8081, "Port to listen on")
	flag.Parse()

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// Main handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Simulate some processing time
		time.Sleep(100 * time.Millisecond)

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Response from backend server on port %d\n", *port)
		fmt.Fprintf(w, "Request: %s %s\n", r.Method, r.URL.Path)
		fmt.Fprintf(w, "Time: %s\n", time.Now().Format(time.RFC3339))
	})

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Backend server starting on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
