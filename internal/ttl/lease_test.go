package ttl_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/ttl"
)

func TestNewLease_NotExpiredImmediately(t *testing.T) {
	l := ttl.NewLease(10 * time.Second)
	if l.Expired() {
		t.Error("lease should not be expired immediately after creation")
	}
}

func TestLease_Remaining_Positive(t *testing.T) {
	l := ttl.NewLease(1 * time.Hour)
	remaining := l.Remaining()
	if remaining <= 0 {
		t.Errorf("expected positive remaining, got %v", remaining)
	}
	if remaining > timett.Errorf("remaining %v exceeds lease duration", remaining)
	}
}

func TestLease_Expired_AfterDuration(t *testing.T) {
	l := ttl.NewLease(10 * time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	if !l.Expired() {
		t.Error("lease should be expired after duration has passed")
	}
	if l.Remaining() != 0 {
		t.Errorf("expected remaining to be 0, got %v", l.Remaining())
	}
}

func TestLease_Renew_ExtendsExpiry(t *testing.T) {
	l := ttl.NewLease(20 * time.Millisecond)
	time.Sleep(15 * time.Millisecond)
	l.Renew()
	time.Sleep(10 * time.Millisecond)
	// After renew, should not yet be expired (only 10ms of new 20ms window elapsed)
	if l.Expired() {
		t.Error("lease should not be expired after renew")
	}
}

func TestLease_ExpiresAt(t *testing.T) {
	before := time.Now()
	d := 5 * time.Minute
	l := ttl.NewLease(d)
	after := time.Now()

	expiry := l.ExpiresAt()
	if expiry.Before(before.Add(d)) {
		t.Errorf("expiry %v is before expected lower bound", expiry)
	}
	if expiry.After(after.Add(d)) {
		t.Errorf("expiry %v is after expected upper bound", expiry)
	}
}
