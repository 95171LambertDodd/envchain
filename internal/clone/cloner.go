// Package clone provides utilities for deep-copying layers and chains
// within envchain, enabling safe branching of environment configurations
// without mutating the originals.
package clone

import (
	"fmt"

	"github.com/yourorg/envchain/internal/config"
)

// Cloner creates independent copies of layers and chains.
type Cloner struct {
	suffix string
}

// NewCloner returns a Cloner that appends suffix to cloned layer names.
// suffix must not be empty.
func NewCloner(suffix string) (*Cloner, error) {
	if suffix == "" {
		return nil, fmt.Errorf("clone: suffix must not be empty")
	}
	return &Cloner{suffix: suffix}, nil
}

// Layer returns a deep copy of src with its name suffixed.
// All key-value pairs are copied; mutations to the clone do not affect src.
func (c *Cloner) Layer(src *config.Layer) (*config.Layer, error) {
	if src == nil {
		return nil, fmt.Errorf("clone: source layer is nil")
	}
	name := src.Name() + c.suffix
	dst, err := config.NewLayer(name)
	if err != nil {
		return nil, fmt.Errorf("clone: create layer %q: %w", name, err)
	}
	for k, v := range src.All() {
		if err := dst.Set(k, v); err != nil {
			return nil, fmt.Errorf("clone: set key %q: %w", k, err)
		}
	}
	return dst, nil
}

// Chain clones every layer in src and assembles them into a new chain
// preserving the original layer order.
func (c *Cloner) Chain(src *config.Chain) (*config.Chain, error) {
	if src == nil {
		return nil, fmt.Errorf("clone: source chain is nil")
	}
	chain := config.NewChain()
	for _, layer := range src.Layers() {
		cloned, err := c.Layer(layer)
		if err != nil {
			return nil, err
		}
		chain.Push(cloned)
	}
	return chain, nil
}
