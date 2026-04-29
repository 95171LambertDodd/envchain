package config

import "fmt"

// Chain holds an ordered stack of Layers. Later layers override earlier ones.
type Chain struct {
	layers []*Layer
}

// NewChain initialises an empty Chain.
func NewChain() *Chain {
	return &Chain{}
}

// Push appends a layer to the top of the chain.
func (c *Chain) Push(l *Layer) {
	c.layers = append(c.layers, l)
}

// Resolve merges all layers from bottom to top and returns a flat map of
// key-value pairs. Values in higher layers win.
func (c *Chain) Resolve() map[string]string {
	merged := make(map[string]string)
	for _, l := range c.layers {
		for k, v := range l.Values {
			merged[k] = v
		}
	}
	return merged
}

// Get returns the resolved value for key, searching from the top layer down.
func (c *Chain) Get(key string) (string, bool) {
	for i := len(c.layers) - 1; i >= 0; i-- {
		if v, ok := c.layers[i].Get(key); ok {
			return v, true
		}
	}
	return "", false
}

// RequireKeys validates that every key in required is present in the resolved
// chain. Returns an error listing all missing keys.
func (c *Chain) RequireKeys(required []string) error {
	resolved := c.Resolve()
	var missing []string
	for _, k := range required {
		if _, ok := resolved[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required config keys: %v", missing)
	}
	return nil
}

// Layers returns the ordered slice of layers (bottom to top).
func (c *Chain) Layers() []*Layer {
	return c.layers
}
