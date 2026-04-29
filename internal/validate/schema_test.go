package validate_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/validate"
)

func TestBuildValidatorFromSchema(t *testing.T) {
	schema := validate.Schema{
		"APP_ENV": {Required: true, Allowed: []string{"dev", "staging", "prod"}},
		"PORT":    {Required: false, Pattern: `^\d{2,5}$`},
	}

	v, err := validate.BuildValidator(schema)
	if err != nil {
		t.Fatalf("unexpected error building validator: %v", err)
	}

	env := map[string]string{
		"APP_ENV": "prod",
		"PORT":    "3000",
	}
	if errs := v.Validate(env); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestBuildValidatorInvalidPattern(t *testing.T) {
	schema := validate.Schema{
		"KEY": {Pattern: `[invalid`},
	}

	_, err := validate.BuildValidator(schema)
	if err == nil {
		t.Fatal("expected error for invalid regex pattern, got nil")
	}
}

func TestBuildValidatorSchemaViolation(t *testing.T) {
	schema := validate.Schema{
		"APP_ENV": {Required: true, Allowed: []string{"dev", "staging", "prod"}},
	}

	v, err := validate.BuildValidator(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Missing required key
	if errs := v.Validate(map[string]string{}); len(errs) == 0 {
		t.Fatal("expected missing-key error, got none")
	}

	// Invalid allowed value
	if errs := v.Validate(map[string]string{"APP_ENV": "local"}); len(errs) == 0 {
		t.Fatal("expected allowed-value error, got none")
	}
}
