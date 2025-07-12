package ingest

import (
	"sync"
	"time"
)

type ConnectionTracker struct {
	mu          sync.RWMutex
	connections map[string]time.Time
	limit       int
}

func NewConnectionTracker(limit int) *ConnectionTracker {
	return &ConnectionTracker{
		connections: make(map[string]time.Time),
		limit:       limit,
	}
}

func (t *ConnectionTracker) Register(id string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.connections) >= t.limit {
		return false
	}
	t.connections[id] = time.Now().UTC()
	return true
}

func (t *ConnectionTracker) Release(id string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.connections, id)
}

func (t *ConnectionTracker) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.connections)
}

func (t *ConnectionTracker) Prune(maxAge time.Duration) int {
	now := time.Now().UTC()
	removed := 0
	t.mu.Lock()
	defer t.mu.Unlock()
	for id, seen := range t.connections {
		if now.Sub(seen) > maxAge {
			delete(t.connections, id)
			removed++
		}
	}
	return removed
}
