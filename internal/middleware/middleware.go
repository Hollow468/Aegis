package middleware

import (
	"apigateway/internal/context"
)

// Chain composes multiple middlewares into a single HandlerFunc.
func Chain(middlewares ...context.HandlerFunc) context.HandlerFunc {
	return func(c *context.GatewayContext) {
		c.SetHandlers(middlewares)
		c.Next()
	}
}
