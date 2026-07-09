// Package main is the entry point for the HTTP server.
// It wires together the router and handlers, then starts
// listening for incoming requests. This is kept minimal —
// it only knows HOW to start the server, not WHAT the
// server does. The actual logic lives in internal/handler.

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/meedaycodes/day01-hello-service/internal/handler"
)

// main initializes the HTTP router, registers route-to-handler
// mappings, and starts the server on port 8080.
// It blocks indefinitely while serving requests.
// If the server fails to start, it logs the error and exits.

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheck)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
