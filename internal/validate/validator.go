package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Key      string
	Required bool
	Pattern  *regexp.Regexp
	Allowed  []string
}

// Validator holds a set of rules and validates resolved configs.
type Validator struct {
	rules []Rule
}

// NewValidator creates a new Validator.
func NewValidator() *Validator {
	return &Validator{}
}

// AddRule appends a validation rule.
func (v *Validator) AddRule(r Rule) {
	v.rules = append(v.rules, r)
}

// Validate checks the provided key-value map against all rules.
// Returns a slice of validation errors (nil if all pass).
func (v *Validator) Validate(env map[string]string) []error {
	var errs []error

	for _, rule := range v.rules {
		val, exists := env[rule.Key]

		if rule.Required && !exists {
			errs = append(errs, fmt.Errorf("required key %q is missing", rule.Key))
			continue
		}

		if !exists {
			continue
		}

		if rule.Pattern != nil && !rule.Pattern.MatchString(val) {
			errs = append(errs, fmt.Errorf("key %q value %q does not match pattern %q", rule.Key, val, rule.Pattern.String()))
		}

		if len(rule.Allowed) > 0 && !contains(rule.Allowed, val) {
			errs = append(errs, fmt.Errorf("key %q value %q is not in allowed set [%s]", rule.Key, val, strings.Join(rule.Allowed, ", ")))
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidateAndSummarize runs Validate and returns a single combined error
// if any validation failures occurred, or nil if all rules pass.
// Useful when callers want a single error rather than a slice.
func (v *Validator) ValidateAndSummarize(env map[string]string) error {
	errs := v.Validate(env)
	if len(errs) == 0 {
		return nil
	}
	messages := make([]string, len(errs))
	for i, err := range errs {
		messages[i] = err.Error()
	}
	return fmt.Errorf("validation failed with %d error(s): %s", len(errs), strings.Join(messages, "; "))
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
