package profile_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/profile"
)

func TestNewProfileEmptyNameError(t *testing.T) {
	_, err := profile.NewProfile("")
	if err == nil {
		t.Fatal("expected error for empty profile name")
	}
}

func TestNewProfileSuccess(t *testing.T) {
	p, err := profile.NewProfile("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "dev" {
		t.Errorf("expected name 'dev', got %q", p.Name)
	}
}

func TestAddLayerDeduplication(t *testing.T) {
	p, _ := profile.NewProfile("staging")
	_ = p.AddLayer("base")
	_ = p.AddLayer("base")
	if len(p.Layers) != 1 {
		t.Errorf("expected 1 layer after dedup, got %d", len(p.Layers))
	}
}

func TestAddLayerEmptyError(t *testing.T) {
	p, _ := profile.NewProfile("prod")
	if err := p.AddLayer(""); err == nil {
		t.Fatal("expected error for empty layer name")
	}
}

func TestRequireKeyDeduplication(t *testing.T) {
	p, _ := profile.NewProfile("prod")
	_ = p.RequireKey("DATABASE_URL")
	_ = p.RequireKey("DATABASE_URL")
	if len(p.RequiredKeys) != 1 {
		t.Errorf("expected 1 required key after dedup, got %d", len(p.RequiredKeys))
	}
}

func TestRequireKeyEmptyError(t *testing.T) {
	p, _ := profile.NewProfile("prod")
	if err := p.RequireKey(""); err == nil {
		t.Fatal("expected error for empty required key")
	}
}

func TestSetAndGetTag(t *testing.T) {
	p, _ := profile.NewProfile("dev")
	_ = p.SetTag("owner", "platform-team")
	v, ok := p.Tag("owner")
	if !ok || v != "platform-team" {
		t.Errorf("expected tag 'platform-team', got %q (ok=%v)", v, ok)
	}
}

func TestSetTagEmptyKeyError(t *testing.T) {
	p, _ := profile.NewProfile("dev")
	if err := p.SetTag("", "value"); err == nil {
		t.Fatal("expected error for empty tag key")
	}
}

func TestValidateNoLayersError(t *testing.T) {
	p, _ := profile.NewProfile("dev")
	if err := p.Validate(); err == nil {
		t.Fatal("expected validation error when no layers defined")
	}
}

func TestValidateWithLayersOK(t *testing.T) {
	p, _ := profile.NewProfile("dev")
	_ = p.AddLayer("base")
	if err := p.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
