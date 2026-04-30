// Package pipeline provides a composable pipeline for processing
// environment config chains through validation, interpolation, and export.
package pipeline

import (
	"fmt"

	"github.com/user/envchain/internal/config"
	"github.com/user/envchain/internal/export"
	"github.com/user/envchain/internal/interpolate"
	"github.com/user/envchain/internal/validate"
)

// Stage represents a processing step applied to a resolved config map.
type Stage func(data map[string]string) (map[string]string, error)

// Pipeline runs a chain through ordered stages: interpolation, validation, export.
type Pipeline struct {
	chain     *config.Chain
	stages    []Stage
	exporter  *export.Exporter
}

// New creates a Pipeline for the given chain and export format.
func New(chain *config.Chain, format string) (*Pipeline, error) {
	ex, err := export.NewExporter(format)
	if err != nil {
		return nil, fmt.Errorf("pipeline: %w", err)
	}
	return &Pipeline{
		chain:    chain,
		exporter: ex,
	}, nil
}

// WithInterpolation adds a variable-expansion stage using the chain as resolver.
func (p *Pipeline) WithInterpolation() *Pipeline {
	p.stages = append(p.stages, func(data map[string]string) (map[string]string, error) {
		resolver := interpolate.MapResolver(data)
		interp := interpolate.NewInterpolator(resolver)
		out := make(map[string]string, len(data))
		for k, v := range data {
			expanded, err := interp.Expand(v)
			if err != nil {
				return nil, fmt.Errorf("interpolation stage: key %q: %w", k, err)
			}
			out[k] = expanded
		}
		return out, nil
	})
	return p
}

// WithValidation adds a validation stage using the provided validator.
func (p *Pipeline) WithValidation(v *validate.Validator) *Pipeline {
	p.stages = append(p.stages, func(data map[string]string) (map[string]string, error) {
		if errs := v.Validate(data); len(errs) > 0 {
			return nil, fmt.Errorf("validation stage: %v", errs)
		}
		return data, nil
	})
	return p
}

// Run resolves the chain, applies all stages, and returns the exported output.
func (p *Pipeline) Run() (string, error) {
	data := p.chain.Resolve()
	for _, stage := range p.stages {
		var err error
		data, err = stage(data)
		if err != nil {
			return "", err
		}
	}
	return p.exporter.Export(data)
}
