// Package diff provides utilities for comparing two environment layers
// or chains and producing a structured change report.
package diff

import (
	"fmt"
	"sort"
)

// ChangeKind describes the type of change between two environments.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// String returns a human-readable representation of the change.
func (c Change) String() string {
	switch c.Kind {
	case Added:
		return fmt.Sprintf("+ %s=%q", c.Key, c.NewValue)
	case Removed:
		return fmt.Sprintf("- %s=%q", c.Key, c.OldValue)
	case Modified:
		return fmt.Sprintf("~ %s: %q -> %q", c.Key, c.OldValue, c.NewValue)
	}
	return ""
}

// Resolver is any type that can resolve a key to a value and list all keys.
type Resolver interface {
	Get(key string) (string, bool)
	Keys() []string
}

// Differ computes the difference between two Resolvers.
type Differ struct{}

// NewDiffer creates a new Differ.
func NewDiffer() *Differ {
	return &Differ{}
}

// Diff computes the ordered list of changes from base to head.
func (d *Differ) Diff(base, head Resolver) []Change {
	seen := make(map[string]struct{})
	var changes []Change

	for _, k := range base.Keys() {
		seen[k] = struct{}{}
		oldVal, _ := base.Get(k)
		newVal, exists := head.Get(k)
		if !exists {
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: oldVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Kind: Modified, OldValue: oldVal, NewValue: newVal})
		}
	}

	for _, k := range head.Keys() {
		if _, ok := seen[k]; ok {
			continue
		}
		newVal, _ := head.Get(k)
		changes = append(changes, Change{Key: k, Kind: Added, NewValue: newVal})
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})
	return changes
}
