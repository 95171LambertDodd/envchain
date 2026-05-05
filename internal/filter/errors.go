package filter

import "errors"

// ErrNoPatternsProvided is returned when a Filter is created with no patterns.
var ErrNoPatternsProvided = errors.New("filter: at least one pattern must be provided")
