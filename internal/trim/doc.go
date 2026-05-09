// Package trim provides whitespace trimming utilities for envchain layers.
//
// # Overview
//
// When environment configuration values are loaded from files, environment
// variables, or user input they often carry accidental leading/trailing
// whitespace. The Trimmer type normalises these values before they are
// consumed by the rest of the pipeline.
//
// # Modes
//
// Three trim modes are available:
//
//   - TrimBoth  – trims whitespace from values only (default for most uses).
//   - TrimKeys  – trims whitespace from key names, leaving values as-is.
//   - TrimAll   – trims both key names and values.
//
// # Usage
//
//	tr, err := trim.NewTrimmer(layer, trim.TrimBoth)
//	if err != nil {
//		log.Fatal(err)
//	}
//	clean := tr.Apply() // map[string]string with trimmed values
//
// For convenience, ApplyToLayer writes trimmed values directly into any
// destination that exposes a Set(key, value string) error method.
package trim
