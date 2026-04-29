package export

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/envchain/envchain/internal/config"
)

// Format represents the output format for exported environment variables.
type Format string

const (
	FormatEnv    Format = "env"
	FormatExport Format = "export"
	FormatDotenv Format = "dotenv"
)

// Exporter writes resolved environment variables from a Chain to an output.
type Exporter struct {
	format Format
}

// NewExporter creates a new Exporter with the given format.
func NewExporter(format Format) (*Exporter, error) {
	switch format {
	case FormatEnv, FormatExport, FormatDotenv:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
}

// Write resolves all keys from the chain and writes them to w.
func (e *Exporter) Write(w io.Writer, chain *config.Chain) error {
	keys := chain.Keys()
	sort.Strings(keys)

	for _, key := range keys {
		val, ok := chain.Get(key)
		if !ok {
			continue
		}
		line, err := e.formatLine(key, val)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return fmt.Errorf("write error for key %q: %w", key, err)
		}
	}
	return nil
}

func (e *Exporter) formatLine(key, val string) (string, error) {
	if strings.ContainsAny(key, "= \t\n") {
		return "", fmt.Errorf("invalid key %q: contains illegal characters", key)
	}
	escaped := strings.ReplaceAll(val, `"`, `\"`)
	switch e.format {
	case FormatExport:
		return fmt.Sprintf(`export %s="%s"`, key, escaped), nil
	case FormatDotenv, FormatEnv:
		return fmt.Sprintf(`%s="%s"`, key, escaped), nil
	default:
		return "", fmt.Errorf("unsupported format: %q", e.format)
	}
}
