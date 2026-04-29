package proxy

import (
	"fmt"
	"net/http/httptest"
	"sync"
	"testing"

	"apigateway/internal/model"
)

func TestRoundRobinBalancer_Next(t *testing.T) {
	b := NewRoundRobinBalancer()
	upstreams := []model.Upstream{
		{Address: "http://localhost:9001", Weight: 1},
		{Address: "http://localhost:9002", Weight: 1},
		{Address: "http://localhost:9003", Weight: 1},
	}

	seen := make(map[string]bool)
	for i := 0; i < 6; i++ {
		u := b.Next(upstreams, nil)
		if u == nil {
			t.Fatal("Next returned nil")
		}
		seen[u.Address] = true
	}
	if len(seen) != 3 {
		t.Errorf("expected 3 unique upstreams, got %d", len(seen))
	}
}

func TestRoundRobinBalancer_Empty(t *testing.T) {
	b := NewRoundRobinBalancer()
	u := b.Next(nil, nil)
	if u != nil {
		t.Error("empty upstreams should return nil")
	}
}

func TestRandomBalancer_Next(t *testing.T) {
	b := NewRandomBalancer()
	upstreams := []model.Upstream{
		{Address: "http://localhost:9001", Weight: 1},
		{Address: "http://localhost:9002", Weight: 1},
		{Address: "http://localhost:9003", Weight: 1},
	}

	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		u := b.Next(upstreams, nil)
		if u == nil {
			t.Fatal("Next returned nil")
		}
		seen[u.Address] = true
	}
	if len(seen) < 2 {
		t.Errorf("expected at least 2 unique upstreams, got %d", len(seen))
	}
}

func TestRandomBalancer_Empty(t *testing.T) {
	b := NewRandomBalancer()
	u := b.Next(nil, nil)
	if u != nil {
		t.Error("empty upstreams should return nil")
	}
}

func TestWeightedRoundRobinBalancer_Basic(t *testing.T) {
	b := NewWeightedRoundRobinBalancer()
	upstreams := []model.Upstream{
		{Address: "http://localhost:9001", Weight: 5},
		{Address: "http://localhost:9002", Weight: 1},
		{Address: "http://localhost:9003", Weight: 1},
	}

	counts := make(map[string]int)
	for i := 0; i < 70; i++ {
		u := b.Next(upstreams, nil)
		if u == nil {
			t.Fatal("Next returned nil")
		}
		counts[u.Address]++
	}

	// 9001 should get ~5x more than others
	if counts["http://localhost:9001"] < counts["http://localhost:9002"] {
		t.Error("weighted upstream should get more requests")
	}
}

func TestWeightedRoundRobinBalancer_EqualWeights(t *testing.T) {
	b := NewWeightedRoundRobinBalancer()
	upstreams := []model.Upstream{
		{Address: "http://localhost:9001", Weight: 1},
		{Address: "http://localhost:9002", Weight: 1},
	}

	counts := make(map[string]int)
	for i := 0; i < 100; i++ {
		u := b.Next(upstreams, nil)
		counts[u.Address]++
	}

	if counts["http://localhost:9001"] != 50 || counts["http://localhost:9002"] != 50 {
		t.Errorf("equal weights should distribute evenly, got %v", counts)
	}
}

func TestConsistentHashBalancer_SameKey(t *testing.T) {
	b := NewConsistentHashBalancer("header:X-User-ID")
	upstreams := []model.Upstream{
		{Address: "http://localhost:9001", Weight: 1},
		{Address: "http://localhost:9002", Weight: 1},
		{Address: "http://localhost:9003", Weight: 1},
	}

	// Same user should always hit same upstream
	r := httptest.NewRequest("GET", "/api/test", nil)
	r.Header.Set("X-User-ID", "user-123")

	first := b.Next(upstreams, r)
	for i := 0; i < 10; i++ {
		r := httptest.NewRequest("GET", "/api/test", nil)
		r.Header.Set("X-User-ID", "user-123")
		next := b.Next(upstreams, r)
		if next.Address != first.Address {
			t.Errorf("same user should hit same upstream, got %s then %s", first.Address, next.Address)
		}
	}
}

func TestConsistentHashBalancer_DifferentKeys(t *testing.T) {
	b := NewConsistentHashBalancer("header:X-User-ID")
	upstreams := []model.Upstream{
		{Address: "http://localhost:9001", Weight: 1},
		{Address: "http://localhost:9002", Weight: 1},
		{Address: "http://localhost:9003", Weight: 1},
	}

	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		r := httptest.NewRequest("GET", "/api/test", nil)
		r.Header.Set("X-User-ID", "user-"+string(rune('A'+i%26)))
		u := b.Next(upstreams, r)
		seen[u.Address] = true
	}
	if len(seen) < 2 {
		t.Errorf("different users should spread across upstreams, got %d unique", len(seen))
	}
}

func TestNewLoadBalancer(t *testing.T) {
	tests := []struct {
		strategy string
		wantType string
	}{
		{"round_robin", "*proxy.RoundRobinBalancer"},
		{"random", "*proxy.RandomBalancer"},
		{"weighted_round_robin", "*proxy.WeightedRoundRobinBalancer"},
		{"consistent_hash", "*proxy.ConsistentHashBalancer"},
		{"", "*proxy.RoundRobinBalancer"},
		{"unknown", "*proxy.RoundRobinBalancer"},
	}

	for _, tt := range tests {
		b := NewLoadBalancer(tt.strategy, "")
		got := typeOf(b)
		if got != tt.wantType {
			t.Errorf("NewLoadBalancer(%q) = %s, want %s", tt.strategy, got, tt.wantType)
		}
	}
}

func typeOf(v interface{}) string {
	if v == nil {
		return "<nil>"
	}
	// Use fmt to get the type
	return fmt.Sprintf("%T", v)
}

func TestConcurrentAccess(t *testing.T) {
	upstreams := []model.Upstream{
		{Address: "http://localhost:9001", Weight: 1},
		{Address: "http://localhost:9002", Weight: 1},
	}

	balancers := []LoadBalancer{
		NewRoundRobinBalancer(),
		NewRandomBalancer(),
		NewWeightedRoundRobinBalancer(),
		NewConsistentHashBalancer("header:X-Test"),
	}

	for _, b := range balancers {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				r := httptest.NewRequest("GET", "/", nil)
				r.Header.Set("X-Test", "test")
				u := b.Next(upstreams, r)
				if u == nil {
					t.Error("concurrent Next returned nil")
				}
			}()
		}
		wg.Wait()
	}
}
