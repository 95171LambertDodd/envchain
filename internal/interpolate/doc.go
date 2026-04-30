// Package interpolate implements variable substitution for envchain config
// values. It supports two reference syntaxes:
//
//   - $VAR           — simple variable reference
//   - ${VAR}         — braced variable reference
//   - ${VAR:-default} — braced reference with a fallback default value
//
// Resolution is performed by a Resolver function, allowing callers to back
// interpolation with any key-value source (a config Layer, OS environment,
// static map, etc.).
//
// Multiple Resolver sources can be composed with ChainResolver, which
// queries each source in order and returns the first match — mirroring the
// layered override semantics used throughout envchain.
//
// Example:
//
//	r := interpolate.ChainResolver(
//		interpolate.MapResolver(layerEntries),
//		interpolate.MapResolver(osEnv),
//	)
//	i := interpolate.NewInterpolator(r)
//	val, err := i.Expand("${DB_HOST:-localhost}:${DB_PORT:-5432}")
package interpolate
