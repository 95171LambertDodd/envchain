package lint_test

import (
	"testing"

	"github.com/envchain/envchain/internal/lint"
)

func TestNewLinterNoRulesError(t *testing.T) {
	_, err := lint.NewLinter()
	if err == nil {
		t.Fatal("expected error for empty rules, got nil")
	}
}

func TestNewLinterSuccess(t *testing.T) {
	l, err := lint.NewLinter(lint.NoEmptyValues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil linter")
	}
}

func TestNoEmptyValuesFindsEmpty(t *testing.T) {
	l, _ := lint.NewLinter(lint.NoEmptyValues)
	findings := l.Lint(map[string]string{"FOO": "", "BAR": "ok"})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "FOO" {
		t.Errorf("expected key FOO, got %q", findings[0].Key)
	}
	if findings[0].Severity != lint.SeverityWarning {
		t.Errorf("expected warning severity, got %q", findings[0].Severity)
	}
}

func TestUppercaseKeysFindsLowercase(t *testing.T) {
	l, _ := lint.NewLinter(lint.UppercaseKeys)
	findings := l.Lint(map[string]string{"foo": "bar", "BAZ": "qux"})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "uppercase-keys" {
		t.Errorf("expected rule uppercase-keys, got %q", findings[0].Rule)
	}
	if findings[0].Severity != lint.SeverityError {
		t.Errorf("expected error severity, got %q", findings[0].Severity)
	}
}

func TestKeyPatternInvalidRegexp(t *testing.T) {
	_, err := lint.KeyPattern("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestKeyPatternMatchPass(t *testing.T) {
	rule, err := lint.KeyPattern(`^[A-Z_]+$`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	l, _ := lint.NewLinter(rule)
	findings := l.Lint(map[string]string{"VALID_KEY": "value"})
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestKeyPatternMatchFail(t *testing.T) {
	rule, _ := lint.KeyPattern(`^[A-Z_]+$`)
	l, _ := lint.NewLinter(rule)
	findings := l.Lint(map[string]string{"bad-key": "value"})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "key-pattern" {
		t.Errorf("expected rule key-pattern, got %q", findings[0].Rule)
	}
}

func TestMultipleRulesCombined(t *testing.T) {
	l, _ := lint.NewLinter(lint.NoEmptyValues, lint.UppercaseKeys)
	entries := map[string]string{
		"foo": "", // triggers both rules
		"BAR": "ok",
	}
	findings := l.Lint(entries)
	if len(findings) != 2 {
		t.Errorf("expected 2 findings, got %d", len(findings))
	}
}
