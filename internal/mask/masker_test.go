package mask_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/mask"
)

func TestIsSensitiveMatchesSubstring(t *testing.T) {
	m := mask.NewMasker("secret", "password", "token")
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"API_TOKEN", true},
		{"APP_SECRET_KEY", true},
		{"DATABASE_URL", false},
		{"PORT", false},
	}
	for _, tc := range cases {
		got := m.IsSensitive(tc.key)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestMaskValueRedactsSensitive(t *testing.T) {
	m := mask.NewMasker("secret")
	got := m.MaskValue("APP_SECRET", "s3cr3t!")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestMaskValuePassesThroughSafe(t *testing.T) {
	m := mask.NewMasker("secret")
	got := m.MaskValue("LOG_LEVEL", "debug")
	if got != "debug" {
		t.Errorf("expected 'debug', got %q", got)
	}
}

func TestMaskMapRedactsOnlySensitiveKeys(t *testing.T) {
	m := mask.NewMasker("password", "token")
	input := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_TOKEN":   "abc123",
		"APP_ENV":     "production",
		"PORT":        "8080",
	}
	result := m.MaskMap(input)

	if result["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_TOKEN"] != "[REDACTED]" {
		t.Errorf("API_TOKEN should be redacted, got %q", result["API_TOKEN"])
	}
	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged, got %q", result["APP_ENV"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged, got %q", result["PORT"])
	}
}

func TestMaskMapDoesNotMutateOriginal(t *testing.T) {
	m := mask.NewMasker("secret")
	input := map[string]string{"MY_SECRET": "topsecret"}
	_ = m.MaskMap(input)
	if input["MY_SECRET"] != "topsecret" {
		t.Error("original map was mutated")
	}
}

func TestNewMaskerNoPatternsNeverMasks(t *testing.T) {
	m := mask.NewMasker()
	if m.IsSensitive("API_SECRET") {
		t.Error("masker with no patterns should never flag a key as sensitive")
	}
}
