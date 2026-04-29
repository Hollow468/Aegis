package middleware

import (
	"net/http"
	"time"

	"apigateway/internal/context"
	"apigateway/internal/logger"

	"go.uber.org/zap"
)

// statusRecorder wraps ResponseWriter to capture the status code.
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// Logger returns a middleware that logs request details.
func Logger() context.HandlerFunc {
	return func(c *context.GatewayContext) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		rec := &statusRecorder{ResponseWriter: c.Writer, statusCode: 200}
		c.Writer = rec

		c.Next()

		latency := time.Since(start)
		logger.Log.Info("request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", rec.statusCode),
			zap.Duration("latency", latency),
			zap.String("remote_addr", c.Request.RemoteAddr),
		)
	}
}
