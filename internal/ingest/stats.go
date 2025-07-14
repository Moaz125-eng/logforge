package ingest

import (
	"net/http"
	"sync/atomic"
)

type Throughput struct {
	batches  atomic.Uint64
	entries  atomic.Uint64
	dropped  atomic.Uint64
}

func NewThroughput() *Throughput {
	return &Throughput{}
}

func (t *Throughput) RecordBatch(count int) {
	t.batches.Add(1)
	t.entries.Add(uint64(count))
}

func (t *Throughput) Drop() {
	t.dropped.Add(1)
}

func (t *Throughput) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"batches":` + itoa(t.batches.Load()) + `,"entries":` + itoa(t.entries.Load()) + `,"dropped":` + itoa(t.dropped.Load()) + `}`))
	}
}
