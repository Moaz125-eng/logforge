package pipeline

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Pool struct {
	workers int
	jobs    chan logentry.Entry
	jobFn   Job
	wg      sync.WaitGroup
	queued  atomic.Uint64
	done    atomic.Uint64
}

func NewPool(workers int, buffer int, job Job) *Pool {
	return &Pool{
		workers: workers,
		jobs:    make(chan logentry.Entry, buffer),
		jobFn:   job,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		w := NewWorker(i, p.jobs, p.jobFn, &p.wg)
		w.Start(ctx)
	}
}

func (p *Pool) Submit(entry logentry.Entry) {
	p.queued.Add(1)
	select {
	case p.jobs <- entry:
	default:
	}
}

func (p *Pool) Close() {
	close(p.jobs)
	p.wg.Wait()
}

func (p *Pool) Stats() (queued, done uint64) {
	return p.queued.Load(), p.done.Load()
}
