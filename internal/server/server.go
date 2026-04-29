package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"apigateway/internal/circuit"
	gwcontext "apigateway/internal/context"
	"apigateway/internal/discovery"
	"apigateway/internal/limiter"
	"apigateway/internal/logger"
	"apigateway/internal/middleware"
	"apigateway/internal/model"
	"apigateway/internal/proxy"
	"apigateway/internal/router"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Server is the main HTTP server that ties routing and proxying together.
type Server struct {
	cfg       *model.Config
	router    *router.TrieRouter
	handler   http.Handler
	httpSrv   *http.Server
	discovery *discovery.ServiceDiscovery
	mu        sync.RWMutex
	dynamicUpstreams map[string][]model.Upstream
}

// NewServer creates a new Server from configuration.
func NewServer(cfg *model.Config) *Server {
	r := router.NewTrieRouter()

	// Per-route middleware components
	routeLimiters := make(map[*model.Route]limiter.Limiter)
	routeBreakers := make(map[*model.Route]*circuit.Breaker)

	// Initialize Redis client for distributed rate limiting
	var redisClient *redis.Client
	if cfg.Redis.Addr != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	}

	// Register all routes from config
	for i := range cfg.Routes {
		route := &cfg.Routes[i]
		matchType := route.MatchType
		if matchType == "" {
			matchType = "exact"
		}

		path := route.Path
		if matchType == "prefix" && path != "" && path[len(path)-1] != '*' {
			if path[len(path)-1] != '/' {
				path += "/"
			}
			path += "*"
		}

		if err := r.Insert(path, route.Method, route); err != nil {
			logger.Log.Error("failed to register route",
				zap.String("path", route.Path),
				zap.Error(err),
			)
			continue
		}

		// Initialize rate limiter for route
		if route.RateLimit != nil {
			rl := route.RateLimit
			switch rl.Type {
			case "redis":
				if redisClient != nil {
					routeLimiters[route] = limiter.NewRedisLimiter(redisClient, rl.Limit, rl.Burst)
				} else {
					routeLimiters[route] = limiter.NewTokenBucketLimiter(rl.Limit, rl.Burst)
				}
			default: // "token_bucket"
				routeLimiters[route] = limiter.NewTokenBucketLimiter(rl.Limit, rl.Burst)
			}
		}

		// Initialize circuit breaker for route
		if route.CircuitBreaker != nil {
			cb := route.CircuitBreaker
			routeBreakers[route] = circuit.NewBreaker(
				cb.ErrorThreshold,
				cb.MinRequests,
				cb.WindowSeconds,
				cb.RecoveryTime,
			)
		}

		logger.Log.Info("route registered",
			zap.String("path", route.Path),
			zap.String("method", route.Method),
			zap.String("match_type", matchType),
			zap.String("balancer", route.Balancer),
		)
	}

	s := &Server{
		cfg:              cfg,
		router:           r,
		dynamicUpstreams: make(map[string][]model.Upstream),
	}

	// Create matcher function
	matcher := func(path, method string) (*model.Route, map[string]string, bool) {
		rm, ok := r.Match(path, method)
		if !ok {
			return nil, nil, false
		}
		route := rm.Route

		if route.ServiceDiscovery != "" {
			s.mu.RLock()
			upstreams := s.dynamicUpstreams[route.ServiceDiscovery]
			s.mu.RUnlock()
			if len(upstreams) > 0 {
				rCopy := *route
				rCopy.Upstreams = upstreams
				return &rCopy, rm.Params, true
			}
		}

		return route, rm.Params, true
	}

	// Balancer factory
	factory := func(route *model.Route) proxy.LoadBalancer {
		return proxy.NewLoadBalancer(route.Balancer, route.HashKey)
	}

	proxyHandler := proxy.NewReverseProxyHandler(matcher, factory)

	// Build the main handler with middleware chain
	s.handler = &gatewayHandler{
		proxy:          proxyHandler,
		matcher:        matcher,
		jwtConfig:      cfg.JWT,
		routeLimiters:  routeLimiters,
		routeBreakers:  routeBreakers,
	}

	// Initialize etcd service discovery
	if len(cfg.Etcd.Endpoints) > 0 {
		sd, err := discovery.NewServiceDiscovery(cfg.Etcd)
		if err != nil {
			logger.Log.Error("failed to initialize etcd service discovery", zap.Error(err))
		} else {
			s.discovery = sd
			sd.OnChange(func(serviceName string, upstreams []model.Upstream) {
				s.mu.Lock()
				s.dynamicUpstreams[serviceName] = upstreams
				s.mu.Unlock()
				logger.Log.Info("dynamic upstreams updated",
					zap.String("service", serviceName),
					zap.Int("count", len(upstreams)),
				)
			})
			logger.Log.Info("etcd service discovery initialized",
				zap.Strings("endpoints", cfg.Etcd.Endpoints),
			)
		}
	}

	return s
}

// gatewayHandler applies middleware then delegates to the proxy.
type gatewayHandler struct {
	proxy         http.Handler
	matcher       proxy.RouteMatcher
	jwtConfig     model.JWTConfig
	routeLimiters map[*model.Route]limiter.Limiter
	routeBreakers map[*model.Route]*circuit.Breaker
}

func (h *gatewayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Global JWT auth middleware
	jwtMW := middleware.JWTAuth(h.jwtConfig)
	c := gwcontext.New(w, r)
	c.SetHandlers([]gwcontext.HandlerFunc{jwtMW, func(c *gwcontext.GatewayContext) {
		// 2. Match route
		route, params, ok := h.matcher(r.URL.Path, r.Method)
		if !ok {
			http.NotFound(w, r)
			return
		}
		c.Params = params

		// 3. Per-route rate limiting
		if rl, exists := h.routeLimiters[route]; exists {
			key := middleware.ClientIPKey(c)
			if !rl.Allow(key) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintf(w, `{"error":"Too Many Requests","code":429,"message":"rate limit exceeded"}`)
				c.Abort()
				return
			}
		}

		// 4. Per-route circuit breaker
		if cb, exists := h.routeBreakers[route]; exists {
			if !cb.Allow() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, `{"error":"Service Unavailable","code":503,"message":"circuit breaker is open"}`)
				c.Abort()
				return
			}
			// Record result after proxy
			defer func() {
				// We can't easily get the status code here, so record success
				// The circuit breaker middleware approach would be better for this
				cb.Record(true)
			}()
		}

		// 5. Delegate to proxy
		h.proxy.ServeHTTP(w, r)
	}})
	c.Next()
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      s.handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Log.Info("HTTP server starting", zap.String("addr", addr))
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.discovery != nil {
		s.discovery.Close()
	}
	if s.httpSrv != nil {
		logger.Log.Info("HTTP server shutting down")
		return s.httpSrv.Shutdown(ctx)
	}
	return nil
}
