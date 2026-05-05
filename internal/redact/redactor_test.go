package redact

import (
	"strings"
	"testing"
)

func TestIsSensitiveMatchesSubstring(t *testing.T) {
	r := NewRedactor([]string{"secret", "token", "password"}, ModeStars)
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"API_TOKEN", true},
		{"MY_SECRET_KEY", true},
		{"DATABASE_URL", false},
		{"APP_ENV", false},
	}
	for _, tc := range cases {
		if got := r.IsSensitive(tc.key); got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestRedactModeStars(t *testing.T) {
	r := NewRedactor([]string{"secret"}, ModeStars)
	if got := r.Redact("supersecretvalue"); got != "********" {
		t.Errorf("expected '********', got %q", got)
	}
}

func TestRedactModePartial(t *testing.T) {
	r := NewRedactor([]string{"token"}, ModePartial)
	got := r.Redact("abcdef")
	if !strings.HasPrefix(got, "ab") {
		t.Errorf("expected partial to start with 'ab', got %q", got)
	}
	if got != "ab****" {
		t.Errorf("expected 'ab****', got %q", got)
	}
}

func TestRedactModePartialShortValue(t *testing.T) {
	r := NewRedactor([]string{"token"}, ModePartial)
	if got := r.Redact("x"); got != "**" {
		t.Errorf("expected '**' for short value, got %q", got)
	}
}

func TestRedactModeHash(t *testing.T) {
	r := NewRedactor([]string{"pass"}, ModeHash)
	got := r.Redact("hunter2")
	if !strings.HasPrefix(got, "[redacted:") {
		t.Errorf("expected hash prefix, got %q", got)
	}
}

func TestRedactMapDoesNotMutateOriginal(t *testing.T) {
	r := NewRedactor([]string{"secret"}, ModeStars)
	orig := map[string]string{
		"APP_SECRET": "mysecret",
		"APP_ENV":    "production",
	}
	out := r.RedactMap(orig)
	if orig["APP_SECRET"] != "mysecret" {
		t.Error("original map was mutated")
	}
	if out["APP_SECRET"] != "********" {
		t.Errorf("expected redacted value, got %q", out["APP_SECRET"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected passthrough value, got %q", out["APP_ENV"])
	}
}

func TestRedactMapPassesThroughSafeKeys(t *testing.T) {
	r := NewRedactor([]string{"secret"}, ModeStars)
	m := map[string]string{"HOST": "localhost", "PORT": "8080"}
	out := r.RedactMap(m)
	for k, v := range m {
		if out[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, out[k])
		}
	}
}
