// Package freeze provides the ability to lock a set of environment keys,
// preventing further modification once frozen.
package freeze

import (
	"errors"
	"fmt"
	"sync"
)

// ErrFrozen is returned when a write is attempted on a frozen key.
var ErrFrozen = errors.New("key is frozen and cannot be modified")

// ErrAlreadyFrozen is returned when Freeze is called on an already-frozen key.
var ErrAlreadyFrozen = errors.New("key is already frozen")

// Freezer tracks which keys are locked and guards a key/value store.
type Freezer struct {
	mu     sync.RWMutex
	values map[string]string
	frozen map[string]struct{}
}

// NewFreezer returns an initialised Freezer.
func NewFreezer() *Freezer {
	return &Freezer{
		values: make(map[string]string),
		frozen: make(map[string]struct{}),
	}
}

// Set stores a value for the given key. Returns ErrFrozen if the key is locked.
func (f *Freezer) Set(key, value string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.frozen[key]; ok {
		return fmt.Errorf("%w: %s", ErrFrozen, key)
	}
	f.values[key] = value
	return nil
}

// Get retrieves the value for the given key and whether it was found.
func (f *Freezer) Get(key string) (string, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	v, ok := f.values[key]
	return v, ok
}

// Freeze marks a key as immutable. Returns ErrAlreadyFrozen if already locked.
func (f *Freezer) Freeze(key string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.frozen[key]; ok {
		return fmt.Errorf("%w: %s", ErrAlreadyFrozen, key)
	}
	f.frozen[key] = struct{}{}
	return nil
}

// IsFrozen reports whether the given key is currently frozen.
func (f *Freezer) IsFrozen(key string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.frozen[key]
	return ok
}

// FrozenKeys returns a sorted snapshot of all frozen key names.
func (f *Freezer) FrozenKeys() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	keys := make([]string, 0, len(f.frozen))
	for k := range f.frozen {
		keys = append(keys, k)
	}
	return keys
}
