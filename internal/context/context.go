package context

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// HandlerFunc is the function signature for request handlers.
type HandlerFunc func(*GatewayContext)

// GatewayContext wraps http.Request/ResponseWriter with routing params and middleware chain support.
type GatewayContext struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Params     map[string]string
	index      int8
	keys       map[string]interface{}
	mu         sync.RWMutex
	statusCode int
	handlers   []HandlerFunc
	aborted    bool
}

// New creates a new GatewayContext.
func New(w http.ResponseWriter, r *http.Request) *GatewayContext {
	return &GatewayContext{
		Writer:  w,
		Request: r,
		Params:  make(map[string]string),
		keys:    make(map[string]interface{}),
		index:   -1,
	}
}

// SetHandlers sets the middleware/handler chain.
func (c *GatewayContext) SetHandlers(handlers []HandlerFunc) {
	c.handlers = handlers
}

// Next executes the next handler in the middleware chain.
func (c *GatewayContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		if c.aborted {
			return
		}
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort stops the middleware chain.
func (c *GatewayContext) Abort() {
	c.aborted = true
}

// IsAborted returns whether the chain was aborted.
func (c *GatewayContext) IsAborted() bool {
	return c.aborted
}

// Set stores a key-value pair in the context (thread-safe).
func (c *GatewayContext) Set(key string, value interface{}) {
	c.mu.Lock()
	c.keys[key] = value
	c.mu.Unlock()
}

// Get retrieves a value by key from the context (thread-safe).
func (c *GatewayContext) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	val, ok := c.keys[key]
	c.mu.RUnlock()
	return val, ok
}

// Status sets the response status code (for chaining).
func (c *GatewayContext) Status(code int) *GatewayContext {
	c.statusCode = code
	return c
}

// Header sets a response header (for chaining).
func (c *GatewayContext) Header(key, value string) *GatewayContext {
	c.Writer.Header().Set(key, value)
	return c
}

// JSON writes a JSON response.
func (c *GatewayContext) JSON(code int, obj interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(code)
	if err := json.NewEncoder(c.Writer).Encode(obj); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, `{"error":"json encode error"}`)
	}
}

// String writes a text response.
func (c *GatewayContext) String(code int, format string, args ...interface{}) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(code)
	fmt.Fprintf(c.Writer, format, args...)
}
