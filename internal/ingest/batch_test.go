package ingest

import (
	"context"
	"testing"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestBatcherFlushesOnSize(t *testing.T) {
	received := 0
	b := NewBatcher(2, time.Second, func(batch Batch) error {
		received += len(batch.Entries)
		return nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b.Start(ctx)
	b.Enqueue(logentry.New("a", "one", logentry.LevelInfo))
	b.Enqueue(logentry.New("a", "two", logentry.LevelInfo))
	time.Sleep(50 * time.Millisecond)
	cancel()
	b.Wait()
	if received < 2 {
		t.Fatalf("expected at least 2 entries flushed")
	}
}
