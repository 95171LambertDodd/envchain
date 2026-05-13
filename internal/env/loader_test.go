package env_test

import (
	"os"
	"testing"

	"github.com/envchain/envchain/internal/env"
)

func setenv(t *testing.T, key, val string) {
	t.Helper()
	t.Setenv(key, val)
}

func TestLoaderEmptyNameError(t *testing.T) {
	l := env.NewLoader()
	_, err := l.Load("")
	if err == nil {
		t.Fatal("expected error for empty layer name")
	}
}

func TestLoaderLoadsAllEnv(t *testing.T) {
	setenv(t, "TESTVAR_FOO", "bar")
	l := env.NewLoader()
	layer, err := l.Load("base")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, err := layer.Get("TESTVAR_FOO")
	if err != nil {
		t.Fatalf("key not found: %v", err)
	}
	if val != "bar" {
		t.Errorf("expected bar, got %q", val)
	}
}

func TestLoaderPrefixFilter(t *testing.T) {
	setenv(t, "APP_HOST", "localhost")
	setenv(t, "APP_PORT", "8080")
	setenv(t, "OTHER_KEY", "ignored")

	l := env.NewLoader(env.WithPrefix("APP_"))
	layer, err := l.Load("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v, err := layer.Get("HOST"); err != nil || v != "localhost" {
		t.Errorf("HOST: got %q, err %v", v, err)
	}
	if v, err := layer.Get("PORT"); err != nil || v != "8080" {
		t.Errorf("PORT: got %q, err %v", v, err)
	}
	if _, err := layer.Get("OTHER_KEY"); err == nil {
		t.Error("OTHER_KEY should not be present after prefix filter")
	}
}

func TestLoaderStrictPassesWhenPresent(t *testing.T) {
	setenv(t, "STRICT_DB_URL", "postgres://localhost/db")

	l := env.NewLoader(
		env.WithPrefix("STRICT_"),
		env.WithStrict("DB_URL"),
	)
	_, err := l.Load("strict")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoaderStrictFailsWhenMissing(t *testing.T) {
	// Ensure the key is absent.
	os.Unsetenv("STRICT_MISSING_KEY")

	l := env.NewLoader(
		env.WithPrefix("STRICT_"),
		env.WithStrict("MISSING_KEY"),
	)
	_, err := l.Load("strict")
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestLoaderNoPrefix(t *testing.T) {
	setenv(t, "NOPREFIX_VAR", "value")

	l := env.NewLoader()
	layer, err := l.Load("raw")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, err := layer.Get("NOPREFIX_VAR"); err != nil || v != "value" {
		t.Errorf("expected value, got %q err %v", v, err)
	}
}
