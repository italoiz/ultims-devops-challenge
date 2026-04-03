package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/italoiz/ultims-devops-challenge/internal/handler"
	"github.com/italoiz/ultims-devops-challenge/internal/metrics"
)

func main() {
	port := envOrDefault("PORT", "8080")
	metricsPort := envOrDefault("METRICS_PORT", "9090")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", withMetrics("health", healthHandler))
	mux.HandleFunc("/readiness", withMetrics("readiness", readinessHandler))
	mux.HandleFunc("/events", withMetrics("events", handler.Events))

	// Metrics on a separate port so it's not exposed through the ingress.
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	go func() {
		addr := ":" + metricsPort
		log.Printf("metrics server listening on %s", addr)
		if err := http.ListenAndServe(addr, metricsMux); err != nil {
			log.Fatalf("metrics server failed: %v", err)
		}
	}()

	addr := ":" + port
	log.Printf("api server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("api server failed: %v", err)
	}
}

// withMetrics wraps a handler to record request count by endpoint and status code.
func withMetrics(endpoint string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rw, r)
		metrics.APIRequests.WithLabelValues(endpoint, strconv.Itoa(rw.status)).Inc()
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

var startTime = time.Now()

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"uptime": fmt.Sprintf("%.0fs", time.Since(startTime).Seconds()),
	})
}

// readinessHandler checks actual dependency availability.
// Right now the API has no external deps (no DB, no queue), so this just
// confirms the server is ready to accept traffic. In production you'd
// ping your datastore and message broker here.
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ready",
	})
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
