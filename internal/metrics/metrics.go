package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	// RequestsTotal counts total HTTP requests by method, path, and status.
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "apigateway",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// RequestDuration tracks request latency in seconds.
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "apigateway",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	// RequestsInFlight tracks currently active requests.
	RequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "apigateway",
			Subsystem: "http",
			Name:      "requests_in_flight",
			Help:      "Number of HTTP requests currently being processed",
		},
	)

	// UpstreamLatency tracks latency to upstream services.
	UpstreamLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "apigateway",
			Subsystem: "upstream",
			Name:      "latency_seconds",
			Help:      "Upstream response latency in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"upstream", "status"},
	)

	// CircuitBreakerState tracks circuit breaker states.
	CircuitBreakerState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "apigateway",
			Subsystem: "circuit",
			Name:      "breaker_state",
			Help:      "Circuit breaker state (0=closed, 1=open, 2=half-open)",
		},
		[]string{"route"},
	)

	// RateLimitHits counts rate-limited requests.
	RateLimitHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "apigateway",
			Subsystem: "ratelimit",
			Name:      "hits_total",
			Help:      "Total number of rate-limited requests",
		},
		[]string{"route"},
	)
)

func init() {
	prometheus.MustRegister(
		RequestsTotal,
		RequestDuration,
		RequestsInFlight,
		UpstreamLatency,
		CircuitBreakerState,
		RateLimitHits,
	)
}

// Handler returns the Prometheus HTTP metrics handler.
func Handler() http.Handler {
	return promhttp.Handler()
}
