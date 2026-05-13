package ttl_test

import (
	"testing"
	"time"

	"github.com/yourorg/envchain/internal/ttl"
)

// staticSource is a minimal in-memory Source for tests.
type staticSource struct {
	data map[string]string
}

func (s *staticSource) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

func (s *staticSource) Keys() []string {
	out := make([]string, 0, len(s.data))
	for k := range s.data {
		out = append(out, k)
	}
	return out
}

func newStore(t *testing.T, data map[string]string, now func() time.Time) *ttl.TTLStore {
	t.Helper()
	store, err := ttl.NewTTLStore(&staticSource{data: data}, now)
	if err != nil {
		t.Fatalf("NewTTLStore: %v", err)
	}
	return store
}

func TestNewTTLStoreNilSourceError(t *testing.T) {
	_, err := ttl.NewTTLStore(nil, time.Now)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestGetBeforeExpiry(t *testing.T) {
	now := time.Unix(1000, 0)
	store := newStore(t, map[string]string{"KEY": "val"}, func() time.Time { return now })

	_ = store.SetExpiry("KEY", now.Add(10*time.Second))

	v, ok := store.Get("KEY")
	if !ok || v != "val" {
		t.Fatalf("expected (val, true), got (%q, %v)", v, ok)
	}
}

func TestGetAfterExpiry(t *testing.T) {
	now := time.Unix(1000, 0)
	store := newStore(t, map[string]string{"KEY": "val"}, func() time.Time { return now })

	_ = store.SetExpiry("KEY", now.Add(-1*time.Second)) // already expired

	_, ok := store.Get("KEY")
	if ok {
		t.Fatal("expected key to be absent after expiry")
	}
}

func TestSetTTLExpires(t *testing.T) {
	var current time.Time = time.Unix(500, 0)
	store := newStore(t, map[string]string{"X": "1"}, func() time.Time { return current })

	_ = store.SetTTL("X", 5*time.Second) // expires at t=505

	current = time.Unix(504, 0)
	if _, ok := store.Get("X"); !ok {
		t.Fatal("key should still be valid at t=504")
	}

	current = time.Unix(505, 0)
	if _, ok := store.Get("X"); ok {
		t.Fatal("key should be expired at t=505")
	}
}

func TestKeysExcludesExpired(t *testing.T) {
	now := time.Unix(1000, 0)
	store := newStore(t, map[string]string{"A": "1", "B": "2"}, func() time.Time { return now })

	_ = store.SetExpiry("A", now.Add(-1*time.Second))

	keys := store.Keys()
	if len(keys) != 1 || keys[0] != "B" {
		t.Fatalf("expected only [B], got %v", keys)
	}
}

func TestIsExpired(t *testing.T) {
	now := time.Unix(1000, 0)
	store := newStore(t, map[string]string{"Z": "v"}, func() time.Time { return now })

	if store.IsExpired("Z") {
		t.Fatal("key should not be expired before SetExpiry")
	}

	_ = store.SetExpiry("Z", now.Add(-time.Second))
	if !store.IsExpired("Z") {
		t.Fatal("key should be expired after past deadline")
	}
}

func TestSetExpiryEmptyKeyError(t *testing.T) {
	store := newStore(t, map[string]string{}, time.Now)
	if err := store.SetExpiry("", time.Now().Add(time.Minute)); err == nil {
		t.Fatal("expected error for empty key")
	}
}
