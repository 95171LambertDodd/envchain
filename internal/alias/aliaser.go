// Package alias provides key aliasing for environment configs,
// allowing one or more alias keys to transparently resolve to a
// canonical key within a layer or chain.
package alias

import (
	"errors"
	"fmt"
)

// Source is anything that can return a value by key.
type Source interface {
	Get(key string) (string, bool)
}

// Aliaser maps alias keys to canonical keys and resolves values
// through a backing Source.
type Aliaser struct {
	source  Source
	aliases map[string]string // alias -> canonical
}

// NewAliaser creates an Aliaser backed by source.
// source must not be nil.
func NewAliaser(source Source) (*Aliaser, error) {
	if source == nil {
		return nil, errors.New("alias: source must not be nil")
	}
	return &Aliaser{
		source:  source,
		aliases: make(map[string]string),
	}, nil
}

// Add registers alias as an alternative name for canonical.
// Both must be non-empty strings and alias must not already be
// registered to a different canonical key.
func (a *Aliaser) Add(alias, canonical string) error {
	if alias == "" {
		return errors.New("alias: alias key must not be empty")
	}
	if canonical == "" {
		return errors.New("alias: canonical key must not be empty")
	}
	if existing, ok := a.aliases[alias]; ok && existing != canonical {
		return fmt.Errorf("alias: %q already mapped to %q", alias, existing)
	}
	a.aliases[alias] = canonical
	return nil
}

// Resolve returns the value for key, checking aliases when the key
// is not found directly in the source. The second return value
// reports whether a value was found.
func (a *Aliaser) Resolve(key string) (string, bool) {
	if v, ok := a.source.Get(key); ok {
		return v, true
	}
	if canonical, ok := a.aliases[key]; ok {
		if v, ok := a.source.Get(canonical); ok {
			return v, true
		}
	}
	return "", false
}

// Canonical returns the canonical key for the given alias, or the
// key itself when no alias is registered.
func (a *Aliaser) Canonical(key string) string {
	if c, ok := a.aliases[key]; ok {
		return c
	}
	return key
}
