// Package coerce provides typed accessors over a string-keyed config source.
//
// Environment variables and layered configs are inherently stringly-typed.
// The coerce package bridges that gap by wrapping any Source (e.g. a
// config.Layer or a merge.Merger result) and exposing String, Int, Bool,
// and Float accessors with fallback defaults.
//
// Example usage:
//
//	c, err := coerce.NewCoercer(layer)
//	if err != nil { ... }
//
//	port, err := c.Int("PORT", 8080)
//	debug, err := c.Bool("DEBUG", false)
//	threshold, err := c.Float("THRESHOLD", 0.95)
//	name := c.String("APP_NAME", "envchain")
//
// All accessors return the supplied fallback when the key is absent, and
// return a descriptive error when the raw value cannot be converted to the
// requested type.
package coerce
