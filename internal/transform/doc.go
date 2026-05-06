// Package transform provides composable transformation functions for
// environment config key-value pairs within the envchain pipeline.
//
// # Overview
//
// A Transformer accepts one or more TransformFunc values and applies them
// in sequence to every entry in a map[string]string. Each function receives
// the current key and value and returns a (possibly modified) key, value,
// and an optional error.
//
// # Built-in transforms
//
//   - UppercaseKeys  – converts all keys to uppercase
//   - TrimSpace      – strips leading/trailing whitespace from values
//   - PrefixKeys     – prepends a fixed string to every key
//   - StripPrefix    – removes a fixed prefix from keys that carry it
//
// # Example
//
//	tr, err := transform.NewTransformer(
//		transform.TrimSpace(),
//		transform.UppercaseKeys(),
//	)
//	if err != nil { ... }
//	out, err := tr.Apply(layer.Entries())
package transform
