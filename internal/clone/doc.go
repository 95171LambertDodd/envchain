// Package clone implements deep-copy semantics for envchain's Layer and Chain
// types.
//
// # Overview
//
// When branching a configuration — for example, to create a staging variant
// from a production chain — it is important that mutations to the new copy do
// not propagate back to the original. The Cloner type handles this safely.
//
// # Usage
//
//	c, err := clone.NewCloner("-clone")
//	if err != nil { /* handle */ }
//
//	// Clone a single layer
//	copied, err := c.Layer(original)
//
//	// Clone an entire chain (all layers)
//	copiedChain, err := c.Chain(originalChain)
//
// # Naming
//
// The suffix supplied to NewCloner is appended to every cloned layer name,
// making it easy to distinguish copies in logs and audit trails.
package clone
