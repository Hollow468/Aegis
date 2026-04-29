package middleware

import (
	"net/http"

	"apigateway/internal/circuit"
	"apigateway/internal/context"
)

// responseStatus wraps ResponseWriter to capture the status code.
type responseStatus struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseStatus) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// CircuitBreaker returns a middleware that enforces circuit breaking.
func CircuitBreaker(breaker *circuit.Breaker) context.HandlerFunc {
	return func(c *context.GatewayContext) {
		if !breaker.Allow() {
			c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"error":   "Service Unavailable",
				"code":    503,
				"message": "circuit breaker is open",
			})
			c.Abort()
			return
		}

		// Wrap writer to capture status
		rec := &responseStatus{ResponseWriter: c.Writer, statusCode: 200}
		c.Writer = rec

		c.Next()

		// Record success/failure based on status code
		success := rec.statusCode < 500
		breaker.Record(success)
	}
}
