package render_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/render"
)

// mapResolver is a simple Resolver backed by a plain map.
type mapResolver map[string]string

func (m mapResolver) Get(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func TestRenderLiteralString(t *testing.T) {
	r := render.NewRenderer(mapResolver{})
	out, err := r.Render("hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", out)
	}
}

func TestRenderEnvFunc(t *testing.T) {
	r := render.NewRenderer(mapResolver{"APP_HOST": "localhost"})
	out, err := r.Render(`{{ env "APP_HOST" }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "localhost" {
		t.Errorf("expected %q, got %q", "localhost", out)
	}
}

func TestRenderEnvFuncMissingKey(t *testing.T) {
	r := render.NewRenderer(mapResolver{})
	_, err := r.Render(`{{ env "MISSING" }}`)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRenderEnvOrFallback(t *testing.T) {
	r := render.NewRenderer(mapResolver{})
	out, err := r.Render(`{{ envOr "MISSING" "default-val" }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "default-val" {
		t.Errorf("expected %q, got %q", "default-val", out)
	}
}

func TestRenderEnvOrPresentKey(t *testing.T) {
	r := render.NewRenderer(mapResolver{"PORT": "8080"})
	out, err := r.Render(`{{ envOr "PORT" "3000" }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "8080" {
		t.Errorf("expected %q, got %q", "8080", out)
	}
}

func TestRenderInvalidTemplate(t *testing.T) {
	r := render.NewRenderer(mapResolver{})
	_, err := r.Render(`{{ unclosed`)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestRenderComposedExpression(t *testing.T) {
	r := render.NewRenderer(mapResolver{"SCHEME": "https", "HOST": "example.com"})
	out, err := r.Render(`{{ env "SCHEME" }}://{{ env "HOST" }}/api`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "https://example.com/api"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}
