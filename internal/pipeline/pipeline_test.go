package pipeline_test

import (
	"strings"
	"testing"

	"github.com/user/envchain/internal/config"
	"github.com/user/envchain/internal/pipeline"
	"github.com/user/envchain/internal/validate"
)

func makeChain(pairs ...string) *config.Chain {
	l := config.NewLayer("test")
	for i := 0; i+1 < len(pairs); i += 2 {
		_ = l.Set(pairs[i], pairs[i+1])
	}
	c := config.NewChain()
	c.Push(l)
	return c
}

func TestPipelineInvalidFormat(t *testing.T) {
	c := makeChain("KEY", "val")
	_, err := pipeline.New(c, "xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestPipelineRunEnvFormat(t *testing.T) {
	c := makeChain("APP_ENV", "production")
	p, err := pipeline.New(c, "env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := p.Run()
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected APP_ENV in output, got: %s", out)
	}
}

func TestPipelineWithInterpolation(t *testing.T) {
	c := makeChain("BASE", "hello", "MSG", "${BASE}_world")
	p, err := pipeline.New(c, "env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p.WithInterpolation()
	out, err := p.Run()
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !strings.Contains(out, "hello_world") {
		t.Errorf("expected interpolated value, got: %s", out)
	}
}

func TestPipelineWithValidationPass(t *testing.T) {
	c := makeChain("PORT", "8080")
	v := validate.NewValidator()
	v.Require("PORT")
	p, err := pipeline.New(c, "env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p.WithValidation(v)
	if _, err := p.Run(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestPipelineWithValidationFail(t *testing.T) {
	c := makeChain("OTHER", "value")
	v := validate.NewValidator()
	v.Require("PORT")
	p, err := pipeline.New(c, "env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p.WithValidation(v)
	_, err = p.Run()
	if err == nil {
		t.Error("expected validation error for missing PORT")
	}
}
