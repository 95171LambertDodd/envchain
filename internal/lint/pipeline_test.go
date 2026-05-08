package lint_test

import (
	"testing"

	"github.com/envchain/envchain/internal/lint"
)

func TestNewLintPipelineNoRulesError(t *testing.T) {
	_, err := lint.NewLintPipeline(nil)
	if err == nil {
		t.Fatal("expected error for no rules")
	}
}

func TestPipelineRunNoFindings(t *testing.T) {
	p, err := lint.NewLintPipeline(nil, lint.NoEmptyValues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rep := p.Run(map[string]string{"FOO": "bar"})
	if rep.HasErrors() {
		t.Errorf("expected no errors, got %d", rep.Errors)
	}
	if len(rep.Findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(rep.Findings))
	}
}

func TestPipelineRunCountsSeverities(t *testing.T) {
	p, err := lint.NewLintPipeline(nil, lint.NoEmptyValues, lint.UppercaseKeys)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "bad" triggers uppercase-keys (error) and no-empty-values (warning)
	rep := p.Run(map[string]string{"bad": ""})
	if rep.Errors != 1 {
		t.Errorf("expected 1 error, got %d", rep.Errors)
	}
	if rep.Warnings != 1 {
		t.Errorf("expected 1 warning, got %d", rep.Warnings)
	}
	if !rep.HasErrors() {
		t.Error("expected HasErrors to be true")
	}
}

func TestPipelineFailOnWarning(t *testing.T) {
	p, err := lint.NewLintPipeline(
		[]lint.PipelineOption{lint.WithFailOnWarning()},
		lint.NoEmptyValues,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rep := p.Run(map[string]string{"FOO": ""})
	// warning promoted to error
	if !rep.HasErrors() {
		t.Error("expected HasErrors true when failOnWarning is set")
	}
	if rep.Warnings != 1 {
		t.Errorf("expected 1 warning recorded, got %d", rep.Warnings)
	}
}

func TestPipelineRunMultipleEntries(t *testing.T) {
	p, err := lint.NewLintPipeline(nil, lint.UppercaseKeys)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rep := p.Run(map[string]string{
		"GOOD":  "v1",
		"bad1":  "v2",
		"Bad2":  "v3",
	})
	if rep.Errors != 2 {
		t.Errorf("expected 2 errors, got %d", rep.Errors)
	}
}
