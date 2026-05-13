// Package chain provides a priority-ordered key-value resolver that fans out
// across multiple [Source] implementations.
//
// # Overview
//
// A [Chainer] holds an ordered slice of sources. When resolving a key the
// Chainer iterates sources from index 0 (highest priority) to the last
// (lowest priority) and returns the value from the first source that contains
// the key. This mirrors the layered override model used throughout envchain.
//
// # Usage
//
//	c, err := chain.New(prodSource, defaultSource)
//	if err != nil { … }
//	_ = c.WithLabel(0, "prod")
//	_ = c.WithLabel(1, "defaults")
//
//	val, ok := c.Get("DATABASE_URL")
//	fmt.Println(c.Origin("DATABASE_URL")) // → "prod"
//
// # Key union
//
// [Chainer.Keys] returns the deduplicated union of all keys present in every
// source, which makes the Chainer itself satisfy the Source interface and
// allows Chainers to be nested.
package chain
