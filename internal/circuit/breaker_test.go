package circuit

import (
	"testing"
	"time"
)

func TestBreaker_Closed_AllowsRequests(t *testing.T) {
	b := NewBreaker(0.5, 10, 60, 30)

	for i := 0; i < 20; i++ {
		if !b.Allow() {
			t.Error("closed breaker should allow requests")
		}
		b.Record(true)
	}
	if b.State() != StateClosed {
		t.Errorf("state = %s, want closed", b.State())
	}
}

func TestBreaker_TripsOnHighErrorRate(t *testing.T) {
	b := NewBreaker(0.5, 10, 60, 30)

	// Send enough requests to meet min_requests
	for i := 0; i < 5; i++ {
		b.Allow()
		b.Record(true)
	}
	for i := 0; i < 5; i++ {
		b.Allow()
		b.Record(false)
	}

	// 50% error rate, exactly at threshold — should trip
	if b.State() != StateOpen {
		t.Errorf("state = %s, want open", b.State())
	}
}

func TestBreaker_Open_RejectsRequests(t *testing.T) {
	b := NewBreaker(0.5, 5, 60, 1) // 1 second recovery

	// Trip the breaker
	for i := 0; i < 3; i++ {
		b.Allow()
		b.Record(true)
	}
	for i := 0; i < 3; i++ {
		b.Allow()
		b.Record(false)
	}

	if b.State() != StateOpen {
		t.Fatalf("state = %s, want open", b.State())
	}

	// Should reject
	if b.Allow() {
		t.Error("open breaker should reject requests")
	}
}

func TestBreaker_HalfOpen_TransitionsToClosed(t *testing.T) {
	b := NewBreaker(0.5, 2, 60, 1) // 1 second recovery

	// Trip the breaker: 2 success + 2 failure = 50% error rate
	for i := 0; i < 2; i++ {
		b.Allow()
		b.Record(true)
	}
	for i := 0; i < 2; i++ {
		b.Allow()
		b.Record(false)
	}

	if b.State() != StateOpen {
		t.Fatalf("state = %s, want open", b.State())
	}

	// Wait for recovery
	time.Sleep(1100 * time.Millisecond)

	// Should transition to half-open
	if !b.Allow() {
		t.Error("half-open breaker should allow test request")
	}

	// Record successes in half-open — should close
	b.Record(true)
	b.Allow()
	b.Record(true)

	if b.State() != StateClosed {
		t.Errorf("state = %s, want closed after successful half-open", b.State())
	}
}

func TestBreaker_HalfOpen_TripsOnFailure(t *testing.T) {
	b := NewBreaker(0.5, 2, 60, 1)

	// Trip the breaker
	for i := 0; i < 2; i++ {
		b.Allow()
		b.Record(true)
	}
	for i := 0; i < 2; i++ {
		b.Allow()
		b.Record(false)
	}

	time.Sleep(1100 * time.Millisecond)

	// Half-open, then fail
	b.Allow()
	b.Record(false)

	if b.State() != StateOpen {
		t.Errorf("state = %s, want open after half-open failure", b.State())
	}
}

func TestBreaker_Reset(t *testing.T) {
	b := NewBreaker(0.5, 3, 60, 30)

	// Trip the breaker
	for i := 0; i < 2; i++ {
		b.Allow()
		b.Record(true)
	}
	for i := 0; i < 2; i++ {
		b.Allow()
		b.Record(false)
	}

	b.Reset()

	if b.State() != StateClosed {
		t.Errorf("state = %s, want closed after reset", b.State())
	}
	if !b.Allow() {
		t.Error("reset breaker should allow requests")
	}
}

func TestBreaker_StateString(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{State(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("State(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}
