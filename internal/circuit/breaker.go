package circuit

import (
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed   State = iota // Normal operation
	StateOpen                  // Rejecting requests
	StateHalfOpen              // Testing recovery
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Breaker implements a Hystrix-style circuit breaker with sliding window.
type Breaker struct {
	mu sync.RWMutex

	// Configuration
	errorThreshold float64       // Error rate threshold (0.0-1.0)
	minRequests    int           // Minimum requests before tripping
	windowSize     time.Duration // Sliding window duration
	recoveryTime   time.Duration // Time in Open before moving to Half-Open

	// State
	state     State
	openSince time.Time

	// Sliding window
	window     []requestRecord
	windowIdx  int
	failures   int
	successes  int
}

type requestRecord struct {
	timestamp time.Time
	success   bool
}

// NewBreaker creates a new circuit breaker.
func NewBreaker(errorThreshold float64, minRequests int, windowSeconds, recoverySeconds int) *Breaker {
	windowSize := time.Duration(windowSeconds) * time.Second
	if windowSize <= 0 {
		windowSize = 60 * time.Second
	}
	recoveryTime := time.Duration(recoverySeconds) * time.Second
	if recoveryTime <= 0 {
		recoveryTime = 30 * time.Second
	}
	if errorThreshold <= 0 || errorThreshold > 1 {
		errorThreshold = 0.5
	}
	if minRequests <= 0 {
		minRequests = 10
	}

	return &Breaker{
		errorThreshold: errorThreshold,
		minRequests:    minRequests,
		windowSize:     windowSize,
		recoveryTime:   recoveryTime,
		state:          StateClosed,
	}
}

// Allow checks if a request should be allowed through.
func (b *Breaker) Allow() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.openSince) >= b.recoveryTime {
			// Transition to Half-Open
			b.mu.RUnlock()
			b.mu.Lock()
			if b.state == StateOpen && time.Since(b.openSince) >= b.recoveryTime {
				b.state = StateHalfOpen
				// Reset window for clean evaluation in Half-Open
				b.window = nil
				b.failures = 0
				b.successes = 0
			}
			b.mu.Unlock()
			b.mu.RLock()
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return true
	}
}

// Record records a request result and updates the circuit state.
func (b *Breaker) Record(success bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()

	// Add to sliding window
	b.window = append(b.window, requestRecord{
		timestamp: now,
		success:   success,
	})

	// Trim old entries from window
	cutoff := now.Add(-b.windowSize)
	start := 0
	for start < len(b.window) && b.window[start].timestamp.Before(cutoff) {
		start++
	}
	if start > 0 {
		b.window = b.window[start:]
	}

	// Recalculate counts
	b.failures = 0
	b.successes = 0
	for _, r := range b.window {
		if r.success {
			b.successes++
		} else {
			b.failures++
		}
	}

	// State transitions
	total := b.failures + b.successes

	switch b.state {
	case StateClosed:
		if total >= b.minRequests {
			errorRate := float64(b.failures) / float64(total)
			if errorRate >= b.errorThreshold {
				b.state = StateOpen
				b.openSince = now
			}
		}

	case StateHalfOpen:
		if !success {
			// Any failure in Half-Open goes back to Open
			b.state = StateOpen
			b.openSince = now
		} else {
			// Success in Half-Open — check if error rate is acceptable
			if total >= b.minRequests {
				errorRate := float64(b.failures) / float64(total)
				if errorRate < b.errorThreshold {
					b.state = StateClosed
				}
			}
		}
	}
}

// State returns the current circuit breaker state.
func (b *Breaker) State() State {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.state
}

// Reset resets the circuit breaker to Closed state.
func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = StateClosed
	b.window = nil
	b.failures = 0
	b.successes = 0
}
