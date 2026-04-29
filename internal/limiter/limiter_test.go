package limiter

import (
	"sync"
	"testing"
	"time"
)

func TestTokenBucketLimiter_Allow(t *testing.T) {
	l := NewTokenBucketLimiter(10, 10) // 10 req/s, burst 10

	// Should allow up to burst
	for i := 0; i < 10; i++ {
		if !l.Allow("client1") {
			t.Errorf("request %d should be allowed", i)
		}
	}

	// 11th request should be denied (burst exhausted)
	if l.Allow("client1") {
		t.Error("11th request should be denied")
	}
}

func TestTokenBucketLimiter_DifferentKeys(t *testing.T) {
	l := NewTokenBucketLimiter(10, 5)

	// Different keys should have separate limits
	for i := 0; i < 5; i++ {
		l.Allow("client1")
	}
	if !l.Allow("client2") {
		t.Error("different client should have separate limit")
	}
}

func TestTokenBucketLimiter_Recovery(t *testing.T) {
	l := NewTokenBucketLimiter(10, 1) // 10/s, burst 1

	// Exhaust burst
	if !l.Allow("client1") {
		t.Fatal("first request should be allowed")
	}
	if l.Allow("client1") {
		t.Error("second request should be denied")
	}

	// Wait for token recovery
	time.Sleep(150 * time.Millisecond)
	if !l.Allow("client1") {
		t.Error("request should be allowed after recovery")
	}
}

func TestTokenBucketLimiter_Concurrent(t *testing.T) {
	l := NewTokenBucketLimiter(100, 100)

	var wg sync.WaitGroup
	allowed := make(chan bool, 200)

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed <- l.Allow("client1")
		}()
	}
	wg.Wait()
	close(allowed)

	count := 0
	for a := range allowed {
		if a {
			count++
		}
	}

	if count > 100 {
		t.Errorf("allowed %d requests, expected <= 100", count)
	}
}

func TestGlobalTokenBucketLimiter(t *testing.T) {
	l := NewGlobalTokenBucketLimiter(10, 5)

	for i := 0; i < 5; i++ {
		if !l.Allow("") {
			t.Errorf("request %d should be allowed", i)
		}
	}
	if l.Allow("") {
		t.Error("6th request should be denied")
	}
}
