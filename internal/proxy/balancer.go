package proxy

import (
	"hash/crc32"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"apigateway/internal/model"
)

// LoadBalancer abstracts the load balancing strategy.
type LoadBalancer interface {
	Next(upstreams []model.Upstream, r *http.Request) *model.Upstream
}

// RoundRobinBalancer selects upstreams in round-robin order.
type RoundRobinBalancer struct {
	counter uint64
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{}
}

func (b *RoundRobinBalancer) Next(upstreams []model.Upstream, _ *http.Request) *model.Upstream {
	if len(upstreams) == 0 {
		return nil
	}
	if len(upstreams) == 1 {
		return &upstreams[0]
	}
	idx := atomic.AddUint64(&b.counter, 1)
	return &upstreams[idx%uint64(len(upstreams))]
}

// RandomBalancer selects upstreams randomly.
type RandomBalancer struct{}

func NewRandomBalancer() *RandomBalancer {
	return &RandomBalancer{}
}

func (b *RandomBalancer) Next(upstreams []model.Upstream, _ *http.Request) *model.Upstream {
	if len(upstreams) == 0 {
		return nil
	}
	if len(upstreams) == 1 {
		return &upstreams[0]
	}
	return &upstreams[rand.Intn(len(upstreams))]
}

// WeightedRoundRobinBalancer selects upstreams based on weight.
type WeightedRoundRobinBalancer struct {
	mu       sync.Mutex
	current  int
	weight   int
	gcd      int
	maxW     int
}

func NewWeightedRoundRobinBalancer() *WeightedRoundRobinBalancer {
	return &WeightedRoundRobinBalancer{}
}

func (b *WeightedRoundRobinBalancer) Next(upstreams []model.Upstream, _ *http.Request) *model.Upstream {
	if len(upstreams) == 0 {
		return nil
	}
	if len(upstreams) == 1 {
		return &upstreams[0]
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// Initialize if needed
	if b.gcd == 0 {
		b.gcd = computeGCD(upstreams)
		b.maxW = maxWeight(upstreams)
		b.weight = b.maxW
		b.current = -1
	}

	for {
		b.current = (b.current + 1) % len(upstreams)
		if b.current == 0 {
			b.weight -= b.gcd
			if b.weight <= 0 {
				b.weight = b.maxW
				if b.weight == 0 {
					b.weight = 1
				}
			}
		}
		if upstreams[b.current].Weight >= b.weight {
			return &upstreams[b.current]
		}
	}
}

func computeGCD(upstreams []model.Upstream) int {
	if len(upstreams) == 0 {
		return 1
	}
	result := upstreams[0].Weight
	if result <= 0 {
		result = 1
	}
	for _, u := range upstreams[1:] {
		w := u.Weight
		if w <= 0 {
			w = 1
		}
		result = gcd(result, w)
	}
	return result
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a <= 0 {
		return 1
	}
	return a
}

func maxWeight(upstreams []model.Upstream) int {
	max := 0
	for _, u := range upstreams {
		if u.Weight > max {
			max = u.Weight
		}
	}
	if max == 0 {
		return 1
	}
	return max
}

// ConsistentHashBalancer uses consistent hashing to map requests to upstreams.
type ConsistentHashBalancer struct {
	hashKey string // "header:X-Name" or "param:name"
}

func NewConsistentHashBalancer(hashKey string) *ConsistentHashBalancer {
	return &ConsistentHashBalancer{hashKey: hashKey}
}

func (b *ConsistentHashBalancer) Next(upstreams []model.Upstream, r *http.Request) *model.Upstream {
	if len(upstreams) == 0 {
		return nil
	}
	if len(upstreams) == 1 {
		return &upstreams[0]
	}

	key := b.extractKey(r)
	if key == "" {
		// Fallback to remote addr
		key = r.RemoteAddr
	}

	hash := crc32.ChecksumIEEE([]byte(key))
	idx := int(hash) % len(upstreams)
	if idx < 0 {
		idx = -idx
	}
	return &upstreams[idx]
}

func (b *ConsistentHashBalancer) extractKey(r *http.Request) string {
	if b.hashKey == "" {
		return ""
	}

	parts := strings.SplitN(b.hashKey, ":", 2)
	if len(parts) != 2 {
		return ""
	}

	switch parts[0] {
	case "header":
		return r.Header.Get(parts[1])
	case "param":
		// Extract from URL path params — this requires the route match params
		// For now, fall back to query param
		return r.URL.Query().Get(parts[1])
	default:
		return ""
	}
}

// NewLoadBalancer creates a LoadBalancer by strategy name.
func NewLoadBalancer(strategy, hashKey string) LoadBalancer {
	switch strategy {
	case "random":
		return NewRandomBalancer()
	case "weighted_round_robin":
		return NewWeightedRoundRobinBalancer()
	case "consistent_hash":
		return NewConsistentHashBalancer(hashKey)
	default: // "round_robin" or empty
		return NewRoundRobinBalancer()
	}
}
