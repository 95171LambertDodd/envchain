// Package scope provides environment scoping for envchain,
// allowing keys to be namespaced by environment (dev, staging, prod)
// and selectively resolved based on the active scope.
package scope

import (
	"errors"
	"fmt"
	"strings"
)

// ErrUnknownScope is returned when an unregistered scope is activated.
var ErrUnknownScope = errors.New("unknown scope")

// Scope represents a named environment tier (e.g. "dev", "staging", "prod").
type Scope struct {
	name    string
	entries map[string]string
}

// NewScope creates a new Scope with the given name.
// Returns an error if the name is empty.
func NewScope(name string) (*Scope, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("scope name must not be empty")
	}
	return &Scope{name: name, entries: make(map[string]string)}, nil
}

// Name returns the scope's name.
func (s *Scope) Name() string { return s.name }

// Set stores a key/value pair in the scope.
func (s *Scope) Set(key, value string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	s.entries[key] = value
	return nil
}

// Get retrieves a value by key. Returns the value and true if found.
func (s *Scope) Get(key string) (string, bool) {
	v, ok := s.entries[key]
	return v, ok
}

// Keys returns all keys defined in this scope.
func (s *Scope) Keys() []string {
	keys := make([]string, 0, len(s.entries))
	for k := range s.entries {
		keys = append(keys, k)
	}
	return keys
}

// ScopeRegistry manages multiple named scopes and tracks the active one.
type ScopeRegistry struct {
	scopes map[string]*Scope
	active string
}

// NewScopeRegistry creates an empty ScopeRegistry.
func NewScopeRegistry() *ScopeRegistry {
	return &ScopeRegistry{scopes: make(map[string]*Scope)}
}

// Register adds a scope to the registry.
func (r *ScopeRegistry) Register(s *Scope) {
	r.scopes[s.Name()] = s
}

// Activate sets the active scope by name.
func (r *ScopeRegistry) Activate(name string) error {
	if _, ok := r.scopes[name]; !ok {
		return fmt.Errorf("%w: %q", ErrUnknownScope, name)
	}
	r.active = name
	return nil
}

// ActiveScope returns the currently active Scope, or nil if none is set.
func (r *ScopeRegistry) ActiveScope() *Scope {
	return r.scopes[r.active]
}

// Resolve looks up a key in the active scope.
// Returns the value and true if found, or empty string and false otherwise.
func (r *ScopeRegistry) Resolve(key string) (string, bool) {
	sc := r.ActiveScope()
	if sc == nil {
		return "", false
	}
	return sc.Get(key)
}
