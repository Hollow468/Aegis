package proxy

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"apigateway/internal/model"
)

// RouteMatcher is a function that finds a route for a given path and method.
// Returns the route, extracted params, and whether a match was found.
type RouteMatcher func(path, method string) (*model.Route, map[string]string, bool)

// BalancerFactory creates a LoadBalancer for a given route.
type BalancerFactory func(route *model.Route) LoadBalancer

// ReverseProxyHandler forwards requests to upstream services.
type ReverseProxyHandler struct {
	matcher   RouteMatcher
	factory   BalancerFactory
	transport *http.Transport
	balancerCache sync.Map // *model.Route -> LoadBalancer
}

// NewReverseProxyHandler creates a new ReverseProxyHandler.
func NewReverseProxyHandler(matcher RouteMatcher, factory BalancerFactory) *ReverseProxyHandler {
	return &ReverseProxyHandler{
		matcher: matcher,
		factory: factory,
		transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ResponseHeaderTimeout: 30 * time.Second,
		},
	}
}

// getBalancer returns or creates a balancer for the route.
func (h *ReverseProxyHandler) getBalancer(route *model.Route) LoadBalancer {
	if cached, ok := h.balancerCache.Load(route); ok {
		return cached.(LoadBalancer)
	}
	balancer := h.factory(route)
	h.balancerCache.Store(route, balancer)
	return balancer
}

// ServeHTTP implements http.Handler.
func (h *ReverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, _, ok := h.matcher(r.URL.Path, r.Method)
	if !ok {
		writeJSONError(w, http.StatusNotFound, "route not found")
		return
	}

	balancer := h.getBalancer(route)
	upstream := balancer.Next(route.Upstreams, r)
	if upstream == nil {
		writeJSONError(w, http.StatusBadGateway, "no upstream available")
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(parseUpstreamURL(upstream.Address))
	proxy.Transport = h.transport

	// Set forwarded headers
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-For", r.RemoteAddr)
		req.Header.Set("X-Real-IP", r.RemoteAddr)
		if r.TLS != nil {
			req.Header.Set("X-Forwarded-Proto", "https")
		} else {
			req.Header.Set("X-Forwarded-Proto", "http")
		}
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		writeJSONError(w, http.StatusBadGateway, "upstream error: "+err.Error())
	}

	proxy.ServeHTTP(w, r)
}

func writeJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   http.StatusText(code),
		"code":    code,
		"message": message,
	})
}

func parseUpstreamURL(address string) *url.URL {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}
	u, err := url.Parse(address)
	if err != nil {
		return &url.URL{Scheme: "http", Host: address}
	}
	return u
}
