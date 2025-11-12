package ingest

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type DeadLetter struct {
	Entry     logentry.Entry `json:"entry"`
	Error     string         `json:"error"`
	Attempts  int            `json:"attempts"`
	QueuedAt  time.Time      `json:"queued_at"`
}

type DLQ struct {
	mu      sync.Mutex
	items   []DeadLetter
	limit   int
	retries int
}

func NewDLQ(limit, retries int) *DLQ {
	return &DLQ{items: make([]DeadLetter, 0, limit), limit: limit, retries: retries}
}

func (q *DLQ) Push(entry logentry.Entry, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	item := DeadLetter{
		Entry: entry, Error: err.Error(), Attempts: 0, QueuedAt: time.Now().UTC(),
	}
	q.items = append(q.items, item)
	if len(q.items) > q.limit {
		q.items = q.items[len(q.items)-q.limit:]
	}
}

func (q *DLQ) List() []DeadLetter {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]DeadLetter, len(q.items))
	copy(out, q.items)
	return out
}

func (q *DLQ) Retry(replay func(logentry.Entry) error) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	kept := make([]DeadLetter, 0, len(q.items))
	ok := 0
	for _, item := range q.items {
		item.Attempts++
		if err := replay(item.Entry); err != nil {
			if item.Attempts < q.retries {
				kept = append(kept, item)
			}
			continue
		}
		ok++
	}
	q.items = kept
	return ok
}

func (q *DLQ) Register(mux *http.ServeMux, replay func(logentry.Entry) error) {
	mux.HandleFunc("/ingest/dlq", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			n := q.Retry(replay)
			_ = json.NewEncoder(w).Encode(map[string]int{"replayed": n})
			return
		}
		_ = json.NewEncoder(w).Encode(q.List())
	})
}
