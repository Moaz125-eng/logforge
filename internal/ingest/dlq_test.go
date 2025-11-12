package ingest

import (
	"errors"
	"testing"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestDLQRetry(t *testing.T) {
	q := NewDLQ(10, 3)
	q.Push(logentry.New("a", "m", logentry.LevelInfo), errors.New("fail"))
	n := q.Retry(func(e logentry.Entry) error { return nil })
	if n != 1 {
		t.Fatalf("expected 1 replay")
	}
}
