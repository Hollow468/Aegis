package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"apigateway/internal/context"
	"apigateway/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTAuth_ValidToken(t *testing.T) {
	cfg := model.JWTConfig{
		Secret:    "test-secret",
		WhiteList: []string{"/health"},
	}

	// Generate valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(cfg.Secret))

	mw := JWTAuth(cfg)
	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	c := context.New(rec, req)

	executed := false
	c.SetHandlers([]context.HandlerFunc{mw, func(c *context.GatewayContext) {
		executed = true
		c.JSON(200, map[string]string{"ok": "true"})
	}})
	c.Next()

	if !executed {
		t.Error("handler should execute with valid token")
	}
	if rec.Code != 200 {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestJWTAuth_MissingHeader(t *testing.T) {
	cfg := model.JWTConfig{Secret: "test-secret"}
	mw := JWTAuth(cfg)

	req := httptest.NewRequest("GET", "/api/users", nil)
	rec := httptest.NewRecorder()
	c := context.New(rec, req)

	c.SetHandlers([]context.HandlerFunc{mw, func(c *context.GatewayContext) {
		t.Error("handler should not execute")
	}})
	c.Next()

	if rec.Code != 401 {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	cfg := model.JWTConfig{Secret: "test-secret"}
	mw := JWTAuth(cfg)

	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()
	c := context.New(rec, req)

	c.SetHandlers([]context.HandlerFunc{mw, func(c *context.GatewayContext) {
		t.Error("handler should not execute")
	}})
	c.Next()

	if rec.Code != 401 {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	cfg := model.JWTConfig{Secret: "test-secret"}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(-time.Hour).Unix(), // expired
	})
	tokenString, _ := token.SignedString([]byte(cfg.Secret))

	mw := JWTAuth(cfg)
	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	c := context.New(rec, req)

	c.SetHandlers([]context.HandlerFunc{mw, func(c *context.GatewayContext) {
		t.Error("handler should not execute")
	}})
	c.Next()

	if rec.Code != 401 {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}

func TestJWTAuth_WhitelistedPath(t *testing.T) {
	cfg := model.JWTConfig{
		Secret:    "test-secret",
		WhiteList: []string{"/health", "/api/public/*"},
	}
	mw := JWTAuth(cfg)

	tests := []string{"/health", "/api/public/anything", "/api/public/"}
	for _, path := range tests {
		req := httptest.NewRequest("GET", path, nil)
		rec := httptest.NewRecorder()
		c := context.New(rec, req)

		executed := false
		c.SetHandlers([]context.HandlerFunc{mw, func(c *context.GatewayContext) {
			executed = true
		}})
		c.Next()

		if !executed {
			t.Errorf("whitelisted path %s should skip auth", path)
		}
	}
}

func TestJWTAuth_WrongSigningMethod(t *testing.T) {
	cfg := model.JWTConfig{Secret: "test-secret"}
	mw := JWTAuth(cfg)

	// Use "none" signing method (not HMAC)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "user123",
	})
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	c := context.New(rec, req)

	c.SetHandlers([]context.HandlerFunc{mw, func(c *context.GatewayContext) {
		t.Error("handler should not execute")
	}})
	c.Next()

	if rec.Code != 401 {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}
