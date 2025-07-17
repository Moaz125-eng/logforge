package pipeline

import (
	"context"
	"sync"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Job func(logentry.Entry) error

type Worker struct {
	id   int
	jobs <-chan logentry.Entry
	run  Job
	wg   *sync.WaitGroup
}

func NewWorker(id int, jobs <-chan logentry.Entry, run Job, wg *sync.WaitGroup) *Worker {
	return &Worker{id: id, jobs: jobs, run: run, wg: wg}
}

func (w *Worker) Start(ctx context.Context) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case entry, ok := <-w.jobs:
				if !ok {
					return
				}
				_ = w.run(entry)
			}
		}
	}()
}
