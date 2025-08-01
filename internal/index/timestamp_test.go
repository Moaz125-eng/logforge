package index

import (
	"testing"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestTimestampRange(t *testing.T) {
	idx := NewTimestampIndex()
	start := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	mid := start.Add(2 * time.Hour)
	end := start.Add(4 * time.Hour)
	idx.Insert(logentry.Entry{Timestamp: start, Message: "a"})
	idx.Insert(logentry.Entry{Timestamp: mid, Message: "b"})
	idx.Insert(logentry.Entry{Timestamp: end, Message: "c"})
	results := idx.Range(start, mid)
	if len(results) != 2 {
		t.Fatalf("expected 2 results got %d", len(results))
	}
}
