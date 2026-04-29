package middleware

import (
	"strconv"
	"time"

	"apigateway/internal/context"
	"apigateway/internal/metrics"
)

// Metrics returns a middleware that records Prometheus metrics.
func Metrics() context.HandlerFunc {
	return func(c *context.GatewayContext) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		metrics.RequestsInFlight.Inc()
		defer metrics.RequestsInFlight.Dec()

		// Wrap writer to capture status
		rec := &responseStatus{ResponseWriter: c.Writer, statusCode: 200}
		c.Writer = rec

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rec.statusCode)

		metrics.RequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.RequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
