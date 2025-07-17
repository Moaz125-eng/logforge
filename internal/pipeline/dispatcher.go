package pipeline

import (
	"context"
	"sync"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Dispatcher struct {
	in       chan logentry.Entry
	out      chan logentry.Entry
	pipeline *Pipeline
	wg       sync.WaitGroup
}

func NewDispatcher(buffer int, pipeline *Pipeline) *Dispatcher {
	return &Dispatcher{
		in:       make(chan logentry.Entry, buffer),
		out:      make(chan logentry.Entry, buffer),
		pipeline: pipeline,
	}
}

func (d *Dispatcher) Start(ctx context.Context, workers int) {
	for i := 0; i < workers; i++ {
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case entry, ok := <-d.in:
					if !ok {
						return
					}
					processed, err := d.pipeline.Run(ctx, entry)
					if err != nil {
						continue
					}
					select {
					case d.out <- processed:
					default:
					}
				}
			}
		}()
	}
}

func (d *Dispatcher) In() chan<- logentry.Entry {
	return d.in
}

func (d *Dispatcher) Out() <-chan logentry.Entry {
	return d.out
}

func (d *Dispatcher) Wait() {
	close(d.in)
	d.wg.Wait()
	close(d.out)
}
