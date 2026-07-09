// Package handler contains HTTP handlers for the service.
// Handlers receive HTTP requests, process them, and write responses.
// This lives in internal/ because it is an implementation detail —
// no external module should import this directly.

package handler

import (
	"encoding/json"
	"net/http"
)

// HealthResponse defines the JSON structure returned by the health endpoint.
// We use a struct over a map because the response shape is known at compile time,
// giving us type safety and explicit JSON field naming via struct tags.

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HealthCheck handles requests to the /health endpoint.
// It returns a 200 OK with a JSON body indicating the service is alive.
// This is typically the first endpoint built in any service because
// load balancers and orchestrators (Kubernetes) poll it to determine
// if the service is healthy and ready to receive traffic.

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status:  "ok",
		Message: "Service is running",
	}

	json.NewEncoder(w).Encode(response)
}
