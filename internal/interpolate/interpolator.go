// Package interpolate provides variable substitution within environment
// config values, supporting ${VAR} and $VAR syntax with optional defaults.
package interpolate

import (
	"fmt"
	"regexp"
	"strings"
)

// Resolver is a function that returns the value for a given key.
type Resolver func(key string) (string, bool)

// Interpolator expands variable references in string values.
type Interpolator struct {
	resolver Resolver
}

// varPattern matches ${VAR}, ${VAR:-default}, and $VAR forms.
var varPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// NewInterpolator creates an Interpolator backed by the given Resolver.
func NewInterpolator(r Resolver) *Interpolator {
	return &Interpolator{resolver: r}
}

// Expand replaces all variable references in s using the Resolver.
// Returns an error if a variable has no value and no default is provided.
func (i *Interpolator) Expand(s string) (string, error) {
	var expandErr error
	result := varPattern.ReplaceAllStringFunc(s, func(match string) string {
		if expandErr != nil {
			return match
		}
		key, defaultVal, hasDefault := parseMatch(match)
		if val, ok := i.resolver(key); ok {
			return val
		}
		if hasDefault {
			return defaultVal
		}
		expandErr = fmt.Errorf("interpolate: unresolved variable %q", key)
		return match
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

// parseMatch extracts the key and optional default from a matched token.
func parseMatch(match string) (key, defaultVal string, hasDefault bool) {
	// $VAR form
	if !strings.HasPrefix(match, "${") {
		return strings.TrimPrefix(match, "$"), "", false
	}
	// ${VAR} or ${VAR:-default}
	inner := match[2 : len(match)-1]
	if idx := strings.Index(inner, ":-"); idx >= 0 {
		return inner[:idx], inner[idx+2:], true
	}
	return inner, "", false
}
