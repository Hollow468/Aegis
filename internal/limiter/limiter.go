package limiter

import (
	"sync"

	"golang.org/x/time/rate"
)

// Limiter is the interface for rate limiting strategies.
type Limiter interface {
	Allow(key string) bool
}

// TokenBucketLimiter implements per-key token bucket rate limiting.
type TokenBucketLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
	rate     rate.Limit
	burst    int
}

// NewTokenBucketLimiter creates a new token bucket limiter.
// r is requests per second, burst is max burst size.
func NewTokenBucketLimiter(r float64, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(r),
		burst:    burst,
	}
}

func (l *TokenBucketLimiter) Allow(key string) bool {
	l.mu.RLock()
	limiter, exists := l.limiters[key]
	l.mu.RUnlock()

	if !exists {
		l.mu.Lock()
		// Double-check after acquiring write lock
		if limiter, exists = l.limiters[key]; !exists {
			limiter = rate.NewLimiter(l.rate, l.burst)
			l.limiters[key] = limiter
		}
		l.mu.Unlock()
	}

	return limiter.Allow()
}

// GlobalTokenBucketLimiter is a single limiter shared across all requests.
type GlobalTokenBucketLimiter struct {
	limiter *rate.Limiter
}

func NewGlobalTokenBucketLimiter(r float64, burst int) *GlobalTokenBucketLimiter {
	return &GlobalTokenBucketLimiter{
		limiter: rate.NewLimiter(rate.Limit(r), burst),
	}
}

func (l *GlobalTokenBucketLimiter) Allow(_ string) bool {
	return l.limiter.Allow()
}
