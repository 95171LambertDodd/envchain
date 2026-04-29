package export_test

import (
	"strings"
	"testing"

	"github.com/envchain/envchain/internal/config"
	"github.com/envchain/envchain/internal/export"
)

func makeChain(t *testing.T, kvs map[string]string) *config.Chain {
	t.Helper()
	layer := config.NewLayer("test")
	for k, v := range kvs {
		if err := layer.Set(k, v); err != nil {
			t.Fatalf("layer.Set(%q, %q): %v", k, v, err)
		}
	}
	chain := config.NewChain(layer)
	return chain
}

func TestNewExporterInvalidFormat(t *testing.T) {
	_, err := export.NewExporter("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExporterEnvFormat(t *testing.T) {
	chain := makeChain(t, map[string]string{"APP_ENV": "production"})
	expr, err := export.NewExporter(export.FormatEnv)
	if err != nil {
		t.Fatalf("NewExporter: %v", err)
	}
	var buf strings.Builder
	if err := expr.Write(&buf, chain); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != `APP_ENV="production"` {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestExporterExportFormat(t *testing.T) {
	chain := makeChain(t, map[string]string{"DB_URL": "postgres://localhost"})
	expr, err := export.NewExporter(export.FormatExport)
	if err != nil {
		t.Fatalf("NewExporter: %v", err)
	}
	var buf strings.Builder
	if err := expr.Write(&buf, chain); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != `export DB_URL="postgres://localhost"` {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestExporterEscapesQuotes(t *testing.T) {
	chain := makeChain(t, map[string]string{"MSG": `say "hello"`})
	expr, _ := export.NewExporter(export.FormatDotenv)
	var buf strings.Builder
	_ = expr.Write(&buf, chain)
	got := strings.TrimSpace(buf.String())
	expected := `MSG="say \"hello\""`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
