// Package defaulter provides a mechanism for applying default values to
// environment config layers. If a key is absent or empty in the target
// layer, the default value is used instead.
package defaulter

import (
	"errors"
	"fmt"
)

// Source is satisfied by any type that can return all its keys and
// retrieve a value by key — compatible with *config.Layer.
type Source interface {
	Keys() []string
	Get(key string) (string, bool)
}

// Target is satisfied by any type that supports reading and writing
// key/value pairs — compatible with *config.Layer.
type Target interface {
	Source
	Set(key, value string) error
}

// Defaulter applies a set of static default values to a Target,
// filling in only keys that are absent or empty.
type Defaulter struct {
	defaults map[string]string
}

// NewDefaulter creates a Defaulter with the provided default map.
// Returns an error if defaults is nil or empty.
func NewDefaulter(defaults map[string]string) (*Defaulter, error) {
	if len(defaults) == 0 {
		return nil, errors.New("defaulter: defaults map must not be empty")
	}
	copy := make(map[string]string, len(defaults))
	for k, v := range defaults {
		if k == "" {
			return nil, errors.New("defaulter: default key must not be empty")
		}
		copy[k] = v
	}
	return &Defaulter{defaults: copy}, nil
}

// Apply writes each default key/value into target only when the key is
// missing or its current value is the empty string. Returns the first
// Set error encountered, if any.
func (d *Defaulter) Apply(target Target) error {
	for k, v := range d.defaults {
		existing, ok := target.Get(k)
		if ok && existing != "" {
			continue
		}
		if err := target.Set(k, v); err != nil {
			return fmt.Errorf("defaulter: setting key %q: %w", k, err)
		}
	}
	return nil
}

// ApplyFromSource reads all keys from src and uses their values as
// defaults when applying to target. Keys already present and non-empty
// in target are left untouched.
func (d *Defaulter) ApplyFromSource(src Source, target Target) error {
	for _, k := range src.Keys() {
		v, _ := src.Get(k)
		existing, ok := target.Get(k)
		if ok && existing != "" {
			continue
		}
		if err := target.Set(k, v); err != nil {
			return fmt.Errorf("defaulter: setting key %q from source: %w", k, err)
		}
	}
	return nil
}
