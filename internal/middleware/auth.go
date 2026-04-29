package middleware

import (
	"net/http"
	"strings"

	"apigateway/internal/context"
	"apigateway/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth returns a middleware that validates JWT tokens.
func JWTAuth(cfg model.JWTConfig) context.HandlerFunc {
	return func(c *context.GatewayContext) {
		// Check whitelist
		if isWhitelisted(c.Request.URL.Path, cfg.WhiteList) {
			c.Next()
			return
		}

		// Extract token from Authorization header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Unauthorized",
				"message": "missing Authorization header",
			})
			c.Abort()
			return
		}

		// Parse "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Unauthorized",
				"message": "invalid Authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Unauthorized",
				"message": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract claims and store in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("claims", claims)
			if sub, ok := claims["sub"].(string); ok {
				c.Set("user_id", sub)
			}
		}

		c.Next()
	}
}

// isWhitelisted checks if a path matches any whitelist pattern.
func isWhitelisted(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if pattern == path {
			return true
		}
		// Support wildcard suffix: /api/public/*
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
	}
	return false
}
