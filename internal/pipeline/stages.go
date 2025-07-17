package pipeline

import (
	"context"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Stage interface {
	Handle(ctx context.Context, entry logentry.Entry) (logentry.Entry, error)
}

type ParseStage struct {
	parse func(string, string) (logentry.Entry, error)
}

func NewParseStage(parse func(string, string) (logentry.Entry, error)) *ParseStage {
	return &ParseStage{parse: parse}
}

func (s *ParseStage) Handle(ctx context.Context, entry logentry.Entry) (logentry.Entry, error) {
	return s.parse("json", entry.Raw)
}

type EnrichStage struct{}

func (s *EnrichStage) Handle(ctx context.Context, entry logentry.Entry) (logentry.Entry, error) {
	if entry.Fields == nil {
		entry.Fields = make(map[string]string)
	}
	entry.Fields["node"] = "local"
	return entry, nil
}

type Pipeline struct {
	stages []Stage
}

func NewPipeline(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

func (p *Pipeline) Run(ctx context.Context, entry logentry.Entry) (logentry.Entry, error) {
	cur := entry
	var err error
	for _, stage := range p.stages {
		cur, err = stage.Handle(ctx, cur)
		if err != nil {
			return logentry.Entry{}, err
		}
	}
	return cur, nil
}
