// Package alias implements key aliasing for envchain environment
// configs.
//
// # Overview
//
// Over the lifetime of a project, environment variable names often
// change — e.g. DB_URL might be renamed to DATABASE_URL. The alias
// package lets you register the old name as an alias for the new
// canonical name so that both continue to resolve correctly without
// duplicating values in your layer files.
//
// # Usage
//
//	// Wrap any Source (e.g. a config.Layer or config.Chain).
//	a, err := alias.NewAliaser(myChain)
//	if err != nil { ... }
//
//	// Register legacy names.
//	_ = a.Add("DB_URL",   "DATABASE_URL")
//	_ = a.Add("PG_HOST",  "DATABASE_HOST")
//
//	// Resolve by either the alias or the canonical key.
//	v, ok := a.Resolve("DB_URL")        // finds DATABASE_URL
//	v, ok  = a.Resolve("DATABASE_URL")  // direct hit
//
// Aliases are one-directional: resolving the canonical key does not
// look up aliases. Each alias may map to exactly one canonical key;
// attempting to re-register an alias to a different canonical key
// returns an error.
package alias
