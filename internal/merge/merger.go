package merge

import (
	"fmt"

	"github.com/envchain/envchain/internal/config"
)

// Strategy defines how conflicting keys are handled during merge.
type Strategy int

const (
	// StrategyOverride replaces existing keys with new values (default).
	StrategyOverride Strategy = iota
	// StrategyKeepFirst retains the first value seen for a key.
	StrategyKeepFirst
	// StrategyError returns an error if a duplicate key is encountered.
	StrategyError
)

// Merger combines multiple layers into a single layer.
type Merger struct {
	strategy Strategy
}

// NewMerger creates a Merger with the given strategy.
func NewMerger(strategy Strategy) *Merger {
	return &Merger{strategy: strategy}
}

// Merge combines the provided layers into a new layer named resultName.
// Layers are processed left to right; later layers take precedence unless
// the strategy dictates otherwise.
func (m *Merger) Merge(resultName string, layers ...*config.Layer) (*config.Layer, error) {
	out := config.NewLayer(resultName)

	for _, layer := range layers {
		for key, value := range layer.All() {
			_, exists := out.Get(key)
			switch m.strategy {
			case StrategyError:
				if exists {
					return nil, fmt.Errorf("merge conflict: duplicate key %q in layer %q", key, layer.Name())
				}
				if err := out.Set(key, value); err != nil {
					return nil, err
				}
			case StrategyKeepFirst:
				if !exists {
					if err := out.Set(key, value); err != nil {
						return nil, err
					}
				}
			default: // StrategyOverride
				if err := out.Set(key, value); err != nil {
					return nil, err
				}
			}
		}
	}

	return out, nil
}
