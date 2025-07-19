package index

import (
	"sort"
	"sync"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type TimestampIndex struct {
	mu      sync.RWMutex
	entries []timedEntry
}

type timedEntry struct {
	ts    time.Time
	entry logentry.Entry
}

func NewTimestampIndex() *TimestampIndex {
	return &TimestampIndex{}
}

func (idx *TimestampIndex) Insert(entry logentry.Entry) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.entries = append(idx.entries, timedEntry{ts: entry.Timestamp, entry: entry})
}

func (idx *TimestampIndex) Range(from, to time.Time) []logentry.Entry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	out := make([]logentry.Entry, 0)
	for _, item := range idx.entries {
		if !item.ts.Before(from) && !item.ts.After(to) {
			out = append(out, item.entry)
		}
	}
	return out
}

func (idx *TimestampIndex) Sort() {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	sort.Slice(idx.entries, func(i, j int) bool {
		return idx.entries[i].ts.Before(idx.entries[j].ts)
	})
}

func (idx *TimestampIndex) Latest(n int) []logentry.Entry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	if n <= 0 || len(idx.entries) == 0 {
		return nil
	}
	start := len(idx.entries) - n
	if start < 0 {
		start = 0
	}
	out := make([]logentry.Entry, 0, n)
	for _, item := range idx.entries[start:] {
		out = append(out, item.entry)
	}
	return out
}
