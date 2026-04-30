package interpolate_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/interpolate"
)

func staticResolver(m map[string]string) interpolate.Resolver {
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}

func TestExpandSimpleVar(t *testing.T) {
	i := interpolate.NewInterpolator(staticResolver(map[string]string{"HOME": "/home/user"}))
	out, err := i.Expand("path=$HOME/bin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "path=/home/user/bin" {
		t.Errorf("got %q", out)
	}
}

func TestExpandBraceVar(t *testing.T) {
	i := interpolate.NewInterpolator(staticResolver(map[string]string{"ENV": "prod"}))
	out, err := i.Expand("mode=${ENV}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "mode=prod" {
		t.Errorf("got %q", out)
	}
}

func TestExpandDefaultUsed(t *testing.T) {
	i := interpolate.NewInterpolator(staticResolver(map[string]string{}))
	out, err := i.Expand("level=${LOG_LEVEL:-info}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "level=info" {
		t.Errorf("got %q", out)
	}
}

func TestExpandDefaultOverriddenByValue(t *testing.T) {
	i := interpolate.NewInterpolator(staticResolver(map[string]string{"LOG_LEVEL": "debug"}))
	out, err := i.Expand("level=${LOG_LEVEL:-info}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "level=debug" {
		t.Errorf("got %q", out)
	}
}

func TestExpandUnresolvedError(t *testing.T) {
	i := interpolate.NewInterpolator(staticResolver(map[string]string{}))
	_, err := i.Expand("host=${DB_HOST}")
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestExpandNoVars(t *testing.T) {
	i := interpolate.NewInterpolator(staticResolver(map[string]string{}))
	out, err := i.Expand("static-value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "static-value" {
		t.Errorf("got %q", out)
	}
}
