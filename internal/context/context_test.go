package context

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestGatewayContext_SetGet(t *testing.T) {
	c := New(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	c.Set("key1", "value1")
	c.Set("key2", 42)

	v, ok := c.Get("key1")
	if !ok || v != "value1" {
		t.Errorf("Get(key1) = %v, %v; want value1, true", v, ok)
	}

	v, ok = c.Get("key2")
	if !ok || v != 42 {
		t.Errorf("Get(key2) = %v, %v; want 42, true", v, ok)
	}

	_, ok = c.Get("nonexistent")
	if ok {
		t.Error("Get(nonexistent) should return false")
	}
}

func TestGatewayContext_JSON(t *testing.T) {
	rec := httptest.NewRecorder()
	c := New(rec, httptest.NewRequest("GET", "/", nil))

	type resp struct {
		Message string `json:"message"`
	}
	c.JSON(200, resp{Message: "hello"})

	if rec.Code != 200 {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var result resp
	json.NewDecoder(rec.Body).Decode(&result)
	if result.Message != "hello" {
		t.Errorf("body message = %q, want hello", result.Message)
	}
}

func TestGatewayContext_String(t *testing.T) {
	rec := httptest.NewRecorder()
	c := New(rec, httptest.NewRequest("GET", "/", nil))

	c.String(201, "hello %s", "world")

	if rec.Code != 201 {
		t.Errorf("status = %d, want 201", rec.Code)
	}
	if body := rec.Body.String(); body != "hello world" {
		t.Errorf("body = %q, want hello world", body)
	}
}

func TestGatewayContext_Header(t *testing.T) {
	rec := httptest.NewRecorder()
	c := New(rec, httptest.NewRequest("GET", "/", nil))

	c.Header("X-Custom", "test-value").JSON(200, nil)

	if v := rec.Header().Get("X-Custom"); v != "test-value" {
		t.Errorf("X-Custom = %q, want test-value", v)
	}
}

func TestGatewayContext_MiddlewareChain(t *testing.T) {
	rec := httptest.NewRecorder()
	c := New(rec, httptest.NewRequest("GET", "/", nil))

	order := []string{}

	handlers := []HandlerFunc{
		func(c *GatewayContext) {
			order = append(order, "before1")
			c.Next()
			order = append(order, "after1")
		},
		func(c *GatewayContext) {
			order = append(order, "handler")
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	expected := []string{"before1", "handler", "after1"}
	if len(order) != len(expected) {
		t.Fatalf("order len = %d, want %d", len(order), len(expected))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("order[%d] = %q, want %q", i, order[i], v)
		}
	}
}

func TestGatewayContext_Abort(t *testing.T) {
	rec := httptest.NewRecorder()
	c := New(rec, httptest.NewRequest("GET", "/", nil))

	executed := false
	handlers := []HandlerFunc{
		func(c *GatewayContext) {
			c.Abort()
		},
		func(c *GatewayContext) {
			executed = true
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	if executed {
		t.Error("second handler should not execute after Abort")
	}
	if !c.IsAborted() {
		t.Error("IsAborted should be true after Abort")
	}
}

func TestGatewayContext_Params(t *testing.T) {
	c := New(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	c.Params["id"] = "42"
	c.Params["name"] = "test"

	if c.Params["id"] != "42" {
		t.Errorf("Params[id] = %q, want 42", c.Params["id"])
	}
}
