package pin

import (
	"fmt"

	"github.com/yourorg/envchain/internal/config"
)

// PinFromLayer creates a Pinner pre-loaded with every key/value from layer.
// This is useful to snapshot a layer as the authoritative baseline and later
// detect drift in a higher-priority layer.
func PinFromLayer(layer *config.Layer) (*Pinner, error) {
	if layer == nil {
		return nil, fmt.Errorf("pin: layer must not be nil")
	}
	p := NewPinner()
	for _, key := range layer.Keys() {
		val, _ := layer.Get(key)
		if err := p.Pin(key, val); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// ValidateChain checks every pinned key against the resolved view of chain.
// It returns all violations found.
func ValidateChain(p *Pinner, chain *config.Chain) []error {
	return p.Validate(chainSource{chain})
}

// chainSource adapts *config.Chain to the Source interface.
type chainSource struct {
	c *config.Chain
}

func (cs chainSource) Get(key string) (string, bool) {
	return cs.c.Resolve(key)
}
