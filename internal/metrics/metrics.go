package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EventsReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_received_total",
			Help: "Total number of tracking events received, labeled by event type.",
		},
		[]string{"type"},
	)

	EventsProcessingDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "events_processing_duration_seconds",
			Help:    "Time spent processing tracking events.",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
	)

	APIRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total API requests, labeled by endpoint and HTTP status code.",
		},
		[]string{"endpoint", "status"},
	)
)
