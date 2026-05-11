// Package index provides forward and reverse key indexing for environment
// layers managed by envchain.
//
// # Overview
//
// An Indexer is built from any Source (a type that exposes Keys and Get)
// and captures a point-in-time snapshot of its entries. It supports:
//
//   - Forward lookup: retrieve a value by key.
//   - Reverse lookup: find all keys that share a given value.
//   - Duplicate detection: identify keys with identical values, which may
//     indicate misconfiguration in layered configs.
//
// # Usage
//
//	idx, err := index.NewIndexer(myLayer)
//	if err != nil { ... }
//
//	v, ok := idx.Get("DATABASE_URL")
//	keys := idx.KeysForValue("postgres://localhost/dev")
//	hasDups, msg := idx.HasDuplicateValues()
package index
