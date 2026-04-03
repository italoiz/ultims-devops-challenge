package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/italoiz/ultims-devops-challenge/internal/metrics"
)

var validEventTypes = map[string]bool{
	"purchase":     true,
	"add_to_cart":  true,
	"page_view":    true,
	"remove_cart":  true,
	"begin_checkout": true,
}

type TrackingEvent struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type EventResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Events(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, EventResponse{
			Status: "error",
			Error:  "method not allowed",
		})
		return
	}

	var event TrackingEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeJSON(w, http.StatusBadRequest, EventResponse{
			Status: "error",
			Error:  fmt.Sprintf("invalid JSON: %v", err),
		})
		return
	}

	if event.Type == "" {
		writeJSON(w, http.StatusBadRequest, EventResponse{
			Status: "error",
			Error:  "missing required field: type",
		})
		return
	}

	if !validEventTypes[event.Type] {
		writeJSON(w, http.StatusBadRequest, EventResponse{
			Status: "error",
			Error:  fmt.Sprintf("unknown event type: %s", event.Type),
		})
		return
	}

	start := time.Now()

	// Simulate event processing. In production this would write to Kafka/NATS/etc.
	time.Sleep(time.Duration(1+time.Now().UnixNano()%5) * time.Millisecond)

	duration := time.Since(start).Seconds()
	metrics.EventsReceived.WithLabelValues(event.Type).Inc()
	metrics.EventsProcessingDuration.Observe(duration)

	writeJSON(w, http.StatusAccepted, EventResponse{
		Status:  "accepted",
		Message: fmt.Sprintf("event %s queued for processing", event.Type),
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
