// Package scope implements environment scoping for envchain.
//
// A Scope represents a named tier in the deployment pipeline — for example
// "dev", "staging", or "prod". Each scope holds its own set of key/value
// environment entries.
//
// A ScopeRegistry manages a collection of scopes and tracks which one is
// currently active. Only the active scope is consulted when resolving keys,
// making it straightforward to switch between environments at runtime without
// mutating the underlying data.
//
// Typical usage:
//
//	reg := scope.NewScopeRegistry()
//
//	dev, _ := scope.NewScope("dev")
//	_ = dev.Set("DB_HOST", "localhost")
//	reg.Register(dev)
//
//	prod, _ := scope.NewScope("prod")
//	_ = prod.Set("DB_HOST", "db.prod.internal")
//	reg.Register(prod)
//
//	_ = reg.Activate("prod")
//	val, ok := reg.Resolve("DB_HOST") // "db.prod.internal", true
package scope
