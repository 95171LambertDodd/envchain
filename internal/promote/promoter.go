// Package promote handles promotion of environment configurations between profiles.
package promote

import (
	"errors"
	"fmt"

	"github.com/your-org/envchain/internal/config"
)

// Strategy defines how conflicts are handled during promotion.
type Strategy int

const (
	// StrategyOverwrite replaces existing keys in the target layer.
	StrategyOverwrite Strategy = iota
	// StrategySkip skips keys that already exist in the target layer.
	StrategySkip
	// StrategyError returns an error if a conflict is detected.
	StrategyError
)

// Result holds the outcome of a promotion operation.
type Result struct {
	Promoted []string
	Skipped  []string
	Errors   []string
}

// Promoter copies entries from a source layer into a target layer.
type Promoter struct {
	strategy Strategy
	keys     []string // if empty, promote all keys
}

// NewPromoter creates a Promoter with the given strategy.
func NewPromoter(strategy Strategy, keys ...string) *Promoter {
	return &Promoter{strategy: strategy, keys: keys}
}

// Promote copies selected (or all) keys from src into dst according to the strategy.
func (p *Promoter) Promote(src, dst *config.Layer) (Result, error) {
	if src == nil || dst == nil {
		return Result{}, errors.New("promote: source and destination layers must not be nil")
	}

	keys := p.keys
	if len(keys) == 0 {
		keys = src.Keys()
	}

	var result Result
	for _, k := range keys {
		v, ok := src.Get(k)
		if !ok {
			result.Errors = append(result.Errors, fmt.Sprintf("key %q not found in source", k))
			continue
		}

		_, exists := dst.Get(k)
		switch {
		case exists && p.strategy == StrategyError:
			return result, fmt.Errorf("promote: conflict on key %q", k)
		case exists && p.strategy == StrategySkip:
			result.Skipped = append(result.Skipped, k)
		default:
			if err := dst.Set(k, v); err != nil {
				return result, fmt.Errorf("promote: failed to set key %q: %w", k, err)
			}
			result.Promoted = append(result.Promoted, k)
		}
	}
	return result, nil
}
