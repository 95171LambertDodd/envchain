package freeze

import (
	"errors"
	"fmt"
)

// Source is the minimal interface required to read keys from a config layer.
type Source interface {
	Keys() []string
	Get(key string) (string, bool)
}

// FreezeFromSource copies all key/value pairs from src into the Freezer and
// immediately freezes each key so they cannot be overridden later.
// Returns a multi-error listing every key that could not be set.
func FreezeFromSource(f *Freezer, src Source) error {
	var errs []error
	for _, k := range src.Keys() {
		v, _ := src.Get(k)
		if err := f.Set(k, v); err != nil {
			errs = append(errs, fmt.Errorf("set %s: %w", k, err))
			continue
		}
		if err := f.Freeze(k); err != nil {
			errs = append(errs, fmt.Errorf("freeze %s: %w", k, err))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
