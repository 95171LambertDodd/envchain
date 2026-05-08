// Package pin provides key pinning: the ability to lock specific keys to
// expected values and detect drift at runtime.
package pin

import (
	"errors"
	"fmt"
	"sync"
)

// ErrPinViolation is returned when a key's value does not match its pinned value.
var ErrPinViolation = errors.New("pin violation")

// Source is anything that can return a value by key.
type Source interface {
	Get(key string) (string, bool)
}

// Pinner stores expected values for keys and validates them against a source.
type Pinner struct {
	mu   sync.RWMutex
	pins map[string]string
}

// NewPinner returns an empty Pinner.
func NewPinner() *Pinner {
	return &Pinner{pins: make(map[string]string)}
}

// Pin records an expected value for key.
func (p *Pinner) Pin(key, value string) error {
	if key == "" {
		return errors.New("pin: key must not be empty")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pins[key] = value
	return nil
}

// Unpin removes a previously pinned key. It is a no-op if the key is not pinned.
func (p *Pinner) Unpin(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.pins, key)
}

// Pinned returns the set of currently pinned keys.
func (p *Pinner) Pinned() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	keys := make([]string, 0, len(p.pins))
	for k := range p.pins {
		keys = append(keys, k)
	}
	return keys
}

// Validate checks every pinned key against src and returns all violations.
func (p *Pinner) Validate(src Source) []error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var errs []error
	for key, expected := range p.pins {
		actual, ok := src.Get(key)
		if !ok {
			errs = append(errs, fmt.Errorf("%w: key %q not found in source", ErrPinViolation, key))
			continue
		}
		if actual != expected {
			errs = append(errs, fmt.Errorf("%w: key %q expected %q got %q", ErrPinViolation, key, expected, actual))
		}
	}
	return errs
}
