package validate_test

import (
	"regexp"
	"testing"

	"github.com/yourorg/envchain/internal/validate"
)

func TestValidatorRequiredKeyPresent(t *testing.T) {
	v := validate.NewValidator()
	v.AddRule(validate.Rule{Key: "APP_ENV", Required: true})

	errs := v.Validate(map[string]string{"APP_ENV": "production"})
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidatorRequiredKeyMissing(t *testing.T) {
	v := validate.NewValidator()
	v.AddRule(validate.Rule{Key: "APP_ENV", Required: true})

	errs := v.Validate(map[string]string{})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidatorPatternMatch(t *testing.T) {
	v := validate.NewValidator()
	v.AddRule(validate.Rule{
		Key:     "PORT",
		Pattern: regexp.MustCompile(`^\d+$`),
	})

	if errs := v.Validate(map[string]string{"PORT": "8080"}); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	if errs := v.Validate(map[string]string{"PORT": "abc"}); len(errs) == 0 {
		t.Fatal("expected pattern error, got none")
	}
}

func TestValidatorAllowedValues(t *testing.T) {
	v := validate.NewValidator()
	v.AddRule(validate.Rule{
		Key:     "LOG_LEVEL",
		Allowed: []string{"debug", "info", "warn", "error"},
	})

	if errs := v.Validate(map[string]string{"LOG_LEVEL": "info"}); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	if errs := v.Validate(map[string]string{"LOG_LEVEL": "verbose"}); len(errs) == 0 {
		t.Fatal("expected allowed-values error, got none")
	}
}

func TestValidatorMultipleErrors(t *testing.T) {
	v := validate.NewValidator()
	v.AddRule(validate.Rule{Key: "DB_URL", Required: true})
	v.AddRule(validate.Rule{Key: "APP_ENV", Required: true})

	errs := v.Validate(map[string]string{})
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}
