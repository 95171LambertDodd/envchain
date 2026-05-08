package lint

import "fmt"

// Report summarises the outcome of a lint run.
type Report struct {
	Findings []Finding
	Errors   int
	Warnings int
	Infos    int
}

// HasErrors returns true when at least one error-severity finding exists.
func (r *Report) HasErrors() bool { return r.Errors > 0 }

// PipelineOption configures a LintPipeline.
type PipelineOption func(*LintPipeline)

// WithFailOnWarning causes the pipeline to treat warnings as errors.
func WithFailOnWarning() PipelineOption {
	return func(p *LintPipeline) { p.failOnWarning = true }
}

// LintPipeline wraps a Linter and produces structured Reports.
type LintPipeline struct {
	linter        *Linter
	failOnWarning bool
}

// NewLintPipeline constructs a LintPipeline from the given rules and options.
func NewLintPipeline(opts []PipelineOption, rules ...Rule) (*LintPipeline, error) {
	l, err := NewLinter(rules...)
	if err != nil {
		return nil, fmt.Errorf("lint pipeline: %w", err)
	}
	p := &LintPipeline{linter: l}
	for _, o := range opts {
		o(p)
	}
	return p, nil
}

// Run lints entries and returns a Report.
func (p *LintPipeline) Run(entries map[string]string) *Report {
	findings := p.linter.Lint(entries)
	rep := &Report{Findings: findings}
	for _, f := range findings {
		switch f.Severity {
		case SeverityError:
			rep.Errors++
		case SeverityWarning:
			rep.Warnings++
			if p.failOnWarning {
				rep.Errors++
			}
		case SeverityInfo:
			rep.Infos++
		}
	}
	return rep
}
