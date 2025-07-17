package pipeline

import (
	"context"

	"github.com/Moaz125-eng/logforge/internal/config"
	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Service struct {
	pool       *Pool
	dispatcher *Dispatcher
}

func NewService(cfg config.Config, parse func(string, string) (logentry.Entry, error), sink func(logentry.Entry) error) *Service {
	pipe := NewPipeline(NewParseStage(parse), &EnrichStage{})
	dispatcher := NewDispatcher(cfg.BatchSize, pipe)
	pool := NewPool(cfg.WorkerCount, cfg.BatchSize*2, func(entry logentry.Entry) error {
		processed, err := pipe.Run(context.Background(), entry)
		if err != nil {
			return err
		}
		return sink(processed)
	})
	return &Service{pool: pool, dispatcher: dispatcher}
}

func (s *Service) Start(ctx context.Context) {
	s.pool.Start(ctx)
	s.dispatcher.Start(ctx, 2)
}

func (s *Service) Process(entry logentry.Entry) {
	s.pool.Submit(entry)
}

func (s *Service) Close() {
	s.pool.Close()
}
