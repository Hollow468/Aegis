package middleware

import (
	"net/http"

	"apigateway/internal/context"
	"apigateway/internal/limiter"
)

// RateLimit returns a middleware that enforces rate limiting.
func RateLimit(l limiter.Limiter, keyFunc func(*context.GatewayContext) string) context.HandlerFunc {
	return func(c *context.GatewayContext) {
		key := keyFunc(c)
		if !l.Allow(key) {
			c.JSON(http.StatusTooManyRequests, map[string]interface{}{
				"error":   "Too Many Requests",
				"code":    429,
				"message": "rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ClientIPKey extracts client IP as the rate limit key.
func ClientIPKey(c *context.GatewayContext) string {
	ip := c.Request.Header.Get("X-Real-IP")
	if ip == "" {
		ip = c.Request.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.Request.RemoteAddr
	}
	return ip
}

// PathKey uses the request path as the rate limit key.
func PathKey(c *context.GatewayContext) string {
	return c.Request.URL.Path
}
