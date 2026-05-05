package promote_test

import (
	"testing"

	"github.com/your-org/envchain/internal/config"
	"github.com/your-org/envchain/internal/promote"
)

func makeLayer(t *testing.T, pairs map[string]string) *config.Layer {
	t.Helper()
	l := config.NewLayer("test")
	for k, v := range pairs {
		if err := l.Set(k, v); err != nil {
			t.Fatalf("makeLayer: Set(%q): %v", k, err)
		}
	}
	return l
}

func TestPromoteOverwriteStrategy(t *testing.T) {
	src := makeLayer(t, map[string]string{"DB_URL": "new", "API_KEY": "key"})
	dst := makeLayer(t, map[string]string{"DB_URL": "old"})

	p := promote.NewPromoter(promote.StrategyOverwrite)
	res, err := p.Promote(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
	v, _ := dst.Get("DB_URL")
	if v != "new" {
		t.Errorf("expected DB_URL=new, got %q", v)
	}
}

func TestPromoteSkipStrategy(t *testing.T) {
	src := makeLayer(t, map[string]string{"DB_URL": "new", "TIMEOUT": "30s"})
	dst := makeLayer(t, map[string]string{"DB_URL": "old"})

	p := promote.NewPromoter(promote.StrategySkip)
	res, err := p.Promote(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "DB_URL" {
		t.Errorf("expected DB_URL skipped, got %v", res.Skipped)
	}
	v, _ := dst.Get("DB_URL")
	if v != "old" {
		t.Errorf("expected DB_URL unchanged, got %q", v)
	}
}

func TestPromoteErrorStrategyConflict(t *testing.T) {
	src := makeLayer(t, map[string]string{"DB_URL": "new"})
	dst := makeLayer(t, map[string]string{"DB_URL": "old"})

	p := promote.NewPromoter(promote.StrategyError)
	_, err := p.Promote(src, dst)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestPromoteSelectedKeys(t *testing.T) {
	src := makeLayer(t, map[string]string{"DB_URL": "db", "SECRET": "s3cr3t", "PORT": "8080"})
	dst := makeLayer(t, map[string]string{})

	p := promote.NewPromoter(promote.StrategyOverwrite, "PORT")
	res, err := p.Promote(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "PORT" {
		t.Errorf("expected only PORT promoted, got %v", res.Promoted)
	}
	if _, ok := dst.Get("SECRET"); ok {
		t.Error("SECRET should not have been promoted")
	}
}

func TestPromoteNilLayerError(t *testing.T) {
	p := promote.NewPromoter(promote.StrategyOverwrite)
	_, err := p.Promote(nil, config.NewLayer("dst"))
	if err == nil {
		t.Error("expected error for nil source")
	}
}
