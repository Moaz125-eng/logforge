package ingest

import (
	"context"

	"github.com/Moaz125-eng/logforge/internal/config"
	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type PipelineSink struct {
	batcher *Batcher
	inner   Sink
}

func NewPipelineSink(cfg config.Config, inner Sink) *PipelineSink {
	ps := &PipelineSink{inner: inner}
	stats := NewThroughput()
	ps.batcher = NewBatcher(cfg.BatchSize, cfg.BatchWindow, func(batch Batch) error {
		stats.RecordBatch(len(batch.Entries))
		for _, entry := range batch.Entries {
			if err := inner(entry); err != nil {
				return err
			}
		}
		return nil
	})
	return ps
}

func (p *PipelineSink) Start(ctx context.Context) {
	p.batcher.Start(ctx)
}

func (p *PipelineSink) Sink(entry logentry.Entry) error {
	p.batcher.Enqueue(entry)
	return nil
}

func (p *PipelineSink) Wait() {
	p.batcher.Wait()
}

func WrapSink(cfg config.Config, inner Sink) *PipelineSink {
	return NewPipelineSink(cfg, inner)
}
