// Package rotate provides key rotation support for envchain layers.
// It allows replacing the value of an existing key with a new value while
// recording the previous value for audit or rollback purposes.
package rotate

import (
	"errors"
	"fmt"
)

// Source is the minimal interface required to read and write key/value pairs.
type Source interface {
	Get(key string) (string, bool)
	Set(key, value string) error
	Keys() []string
}

// Record captures a single rotation event.
type Record struct {
	Key      string
	OldValue string
	NewValue string
}

// Rotator rotates key values in a Source, optionally restricted to a set of
// allowed keys.
type Rotator struct {
	src     Source
	allowed map[string]struct{}
}

// NewRotator creates a Rotator backed by src. If allowedKeys is non-empty only
// those keys may be rotated; pass nil or an empty slice to allow all keys.
func NewRotator(src Source, allowedKeys []string) (*Rotator, error) {
	if src == nil {
		return nil, errors.New("rotate: source must not be nil")
	}
	allowed := make(map[string]struct{}, len(allowedKeys))
	for _, k := range allowedKeys {
		if k == "" {
			return nil, errors.New("rotate: allowed key must not be empty")
		}
		allowed[k] = struct{}{}
	}
	return &Rotator{src: src, allowed: allowed}, nil
}

// Rotate replaces the value of key with newValue. It returns a Record containing
// the previous value. An error is returned if the key is not present in the
// source, or if the key is not in the allowed set (when one is configured).
func (r *Rotator) Rotate(key, newValue string) (Record, error) {
	if key == "" {
		return Record{}, errors.New("rotate: key must not be empty")
	}
	if len(r.allowed) > 0 {
		if _, ok := r.allowed[key]; !ok {
			return Record{}, fmt.Errorf("rotate: key %q is not in the allowed set", key)
		}
	}
	old, ok := r.src.Get(key)
	if !ok {
		return Record{}, fmt.Errorf("rotate: key %q not found in source", key)
	}
	if err := r.src.Set(key, newValue); err != nil {
		return Record{}, fmt.Errorf("rotate: failed to set key %q: %w", key, err)
	}
	return Record{Key: key, OldValue: old, NewValue: newValue}, nil
}

// RotateAll rotates every key returned by src.Keys() using the values supplied
// in updates. Keys present in updates but absent from the source are skipped
// and reported as errors in the returned slice; rotation continues for
// remaining keys.
func (r *Rotator) RotateAll(updates map[string]string) ([]Record, []error) {
	var records []Record
	var errs []error
	for k, v := range updates {
		rec, err := r.Rotate(k, v)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		records = append(records, rec)
	}
	return records, errs
}
