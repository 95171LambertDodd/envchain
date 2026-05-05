package filter_test

import (
	"testing"

	"envchain/internal/filter"
)

func TestNewFilterNoPatternsError(t *testing.T) {
	_, err := filter.NewFilter(filter.ModePrefix)
	if err == nil {
		t.Fatal("expected error for no patterns, got nil")
	}
}

func TestFilterPrefixApply(t *testing.T) {
	f, _ := filter.NewFilter(filter.ModePrefix, "APP_", "DB_")
	env := map[string]string{
		"APP_HOST": "localhost",
		"DB_PASS":  "secret",
		"LOG_LEVEL": "info",
	}
	result := f.Apply(env)
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := result["DB_PASS"]; !ok {
		t.Error("expected DB_PASS in result")
	}
	if _, ok := result["LOG_LEVEL"]; ok {
		t.Error("did not expect LOG_LEVEL in result")
	}
}

func TestFilterSuffixApply(t *testing.T) {
	f, _ := filter.NewFilter(filter.ModeSuffix, "_URL", "_KEY")
	env := map[string]string{
		"DATABASE_URL": "postgres://...",
		"API_KEY":      "abc123",
		"DEBUG":        "true",
	}
	result := f.Apply(env)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilterExactApply(t *testing.T) {
	f, _ := filter.NewFilter(filter.ModeExact, "PORT", "HOST")
	env := map[string]string{
		"PORT":    "8080",
		"HOST":    "0.0.0.0",
		"TIMEOUT": "30s",
	}
	result := f.Apply(env)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilterExclude(t *testing.T) {
	f, _ := filter.NewFilter(filter.ModePrefix, "SECRET_")
	env := map[string]string{
		"SECRET_TOKEN": "hidden",
		"APP_NAME":     "envchain",
		"VERSION":      "1.0",
	}
	result := f.Exclude(env)
	if _, ok := result["SECRET_TOKEN"]; ok {
		t.Error("did not expect SECRET_TOKEN in excluded result")
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries after exclusion, got %d", len(result))
	}
}

func TestFilterDoesNotMutateInput(t *testing.T) {
	f, _ := filter.NewFilter(filter.ModeExact, "KEY")
	env := map[string]string{"KEY": "val", "OTHER": "x"}
	f.Apply(env)
	if len(env) != 2 {
		t.Error("Apply must not mutate the input map")
	}
}
