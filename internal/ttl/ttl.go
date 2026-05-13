// Package ttl provides time-to-live expiry tracking for environment config entries.
// Keys can be assigned an expiration deadline; expired keys are treated as absent.
package ttl

import (
	"errors"
	"sync"
	"time"
)

// Source is the interface expected from a key/value provider.
type Source interface {
	Keys() []string
	Get(key string) (string, bool)
}

// TTLStore wraps a Source and enforces per-key expiry deadlines.
type TTLStore struct {
	mu       sync.RWMutex
	source   Source
	expiries map[string]time.Time
	now      func() time.Time
}

// NewTTLStore creates a TTLStore backed by source.
// now is injectable for testing; pass time.Now for production use.
func NewTTLStore(source Source, now func() time.Time) (*TTLStore, error) {
	if source == nil {
		return nil, errors.New("ttl: source must not be nil")
	}
	if now == nil {
		now = time.Now
	}
	return &TTLStore{
		source:   source,
		expiries: make(map[string]time.Time),
		now:      now,
	}, nil
}

// SetExpiry assigns an absolute expiry time to key.
// After that instant, Get will report the key as absent.
func (t *TTLStore) SetExpiry(key string, at time.Time) error {
	if key == "" {
		return errors.New("ttl: key must not be empty")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.expiries[key] = at
	return nil
}

// SetTTL is a convenience wrapper that sets expiry to now+d.
func (t *TTLStore) SetTTL(key string, d time.Duration) error {
	return t.SetExpiry(key, t.now().Add(d))
}

// Get returns the value for key from the underlying source, unless the key
// has expired, in which case it returns ("", false).
func (t *TTLStore) Get(key string) (string, bool) {
	t.mu.RLock()
	expiry, hasExpiry := t.expiries[key]
	t.mu.RUnlock()

	if hasExpiry && !t.now().Before(expiry) {
		return "", false
	}
	return t.source.Get(key)
}

// Keys returns all non-expired keys present in the source.
func (t *TTLStore) Keys() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	now := t.now()
	var out []string
	for _, k := range t.source.Keys() {
		if exp, ok := t.expiries[k]; ok && !now.Before(exp) {
			continue
		}
		out = append(out, k)
	}
	return out
}

// IsExpired reports whether key has a recorded expiry that has already passed.
func (t *TTLStore) IsExpired(key string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	exp, ok := t.expiries[key]
	return ok && !t.now().Before(exp)
}
