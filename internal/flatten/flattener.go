// Package flatten provides utilities for collapsing a layered config chain
// into a single flat map, with configurable conflict resolution and prefix support.
package flatten

import (
	"errors"
	"fmt"
)

// Strategy controls how key conflicts are handled during flattening.
type Strategy int

const (
	// StrategyLastWins uses the value from the last layer that defines a key.
	StrategyLastWins Strategy = iota
	// StrategyFirstWins keeps the value from the first layer that defines a key.
	StrategyFirstWins
	// StrategyError returns an error if any key is defined in more than one layer.
	StrategyError
)

// ErrConflict is returned when StrategyError is used and a duplicate key is found.
var ErrConflict = errors.New("flatten: key conflict")

// Source is the minimal interface a config layer must satisfy to be flattened.
type Source interface {
	Name() string
	Keys() []string
	Get(key string) (string, bool)
}

// Flattener collapses multiple Sources into a single map.
type Flattener struct {
	strategy Strategy
	prefix   string
}

// NewFlattener constructs a Flattener with the given strategy.
func NewFlattener(strategy Strategy) *Flattener {
	return &Flattener{strategy: strategy}
}

// WithPrefix returns a new Flattener that prepends prefix to every key.
func (f *Flattener) WithPrefix(prefix string) *Flattener {
	return &Flattener{strategy: f.strategy, prefix: prefix}
}

// Flatten merges sources (ordered lowest-to-highest priority) into a flat map.
func (f *Flattener) Flatten(sources []Source) (map[string]string, error) {
	out := make(map[string]string)
	seen := make(map[string]string) // key -> layer name that first set it

	for _, src := range sources {
		for _, k := range src.Keys() {
			v, _ := src.Get(k)
			outKey := f.prefix + k

			switch f.strategy {
			case StrategyFirstWins:
				if _, exists := out[outKey]; !exists {
					out[outKey] = v
					seen[outKey] = src.Name()
				}
			case StrategyError:
				if prev, exists := seen[outKey]; exists {
					return nil, fmt.Errorf("%w: key %q defined in %q and %q",
						ErrConflict, outKey, prev, src.Name())
				}
				out[outKey] = v
				seen[outKey] = src.Name()
			default: // StrategyLastWins
				out[outKey] = v
				seen[outKey] = src.Name()
			}
		}
	}
	return out, nil
}
