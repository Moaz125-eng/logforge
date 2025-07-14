package ingest

import (
	"context"
	"sync"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Batch struct {
	Entries []logentry.Entry
	Flushed time.Time
}

type Batcher struct {
	size    int
	window  time.Duration
	in      chan logentry.Entry
	out     chan Batch
	sink    func(Batch) error
	wg      sync.WaitGroup
	closed  chan struct{}
}

func NewBatcher(size int, window time.Duration, sink func(Batch) error) *Batcher {
	return &Batcher{
		size:   size,
		window: window,
		in:     make(chan logentry.Entry, size*2),
		out:    make(chan Batch, 8),
		sink:   sink,
		closed: make(chan struct{}),
	}
}

func (b *Batcher) Start(ctx context.Context) {
	b.wg.Add(1)
	go b.loop(ctx)
}

func (b *Batcher) loop(ctx context.Context) {
	defer b.wg.Done()
	buf := make([]logentry.Entry, 0, b.size)
	timer := time.NewTimer(b.window)
	defer timer.Stop()
	flush := func() {
		if len(buf) == 0 {
			return
		}
		batch := Batch{Entries: append([]logentry.Entry(nil), buf...), Flushed: time.Now().UTC()}
		buf = buf[:0]
		if b.sink != nil {
			_ = b.sink(batch)
		}
		select {
		case b.out <- batch:
		default:
		}
	}
	for {
		select {
		case <-ctx.Done():
			flush()
			close(b.closed)
			return
		case entry := <-b.in:
			buf = append(buf, entry)
			if len(buf) >= b.size {
				flush()
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(b.window)
			}
		case <-timer.C:
			flush()
			timer.Reset(b.window)
		}
	}
}

func (b *Batcher) Enqueue(entry logentry.Entry) {
	select {
	case b.in <- entry:
	default:
	}
}

func (b *Batcher) Wait() {
	b.wg.Wait()
	<-b.closed
}

func (b *Batcher) Out() <-chan Batch {
	return b.out
}
