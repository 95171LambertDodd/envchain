// Package rotate provides key rotation utilities for envchain.
//
// A Rotator wraps any Source (a layer, scope, or other key/value store) and
// allows callers to atomically replace the value of an existing key while
// capturing the previous value in a Record.
//
// Basic usage:
//
//	r, err := rotate.NewRotator(layer, nil) // nil = allow all keys
//	if err != nil {
//		log.Fatal(err)
//	}
//	rec, err := r.Rotate("DB_PASSWORD", newSecret)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("rotated %s: %s -> %s\n", rec.Key, rec.OldValue, rec.NewValue)
//
// To restrict which keys may be rotated supply an explicit allow-list:
//
//	r, err := rotate.NewRotator(layer, []string{"DB_PASSWORD", "API_KEY"})
//
// RotateAll accepts a map of updates and processes each entry, collecting
// partial errors without aborting the entire batch.
package rotate
