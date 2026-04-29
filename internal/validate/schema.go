package validate

import (
	"fmt"
	"regexp"
)

// SchemaEntry describes a single field in a schema definition.
type SchemaEntry struct {
	Required bool
	Pattern  string
	Allowed  []string
}

// Schema maps key names to their SchemaEntry definitions.
type Schema map[string]SchemaEntry

// BuildValidator constructs a Validator from a Schema, compiling regex patterns.
func BuildValidator(schema Schema) (*Validator, error) {
	v := NewValidator()

	for key, entry := range schema {
		rule := Rule{
			Key:      key,
			Required: entry.Required,
			Allowed:  entry.Allowed,
		}

		if entry.Pattern != "" {
			re, err := regexp.Compile(entry.Pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid pattern for key %q: %w", key, err)
			}
			rule.Pattern = re
		}

		v.AddRule(rule)
	}

	return v, nil
}
