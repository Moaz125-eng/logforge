package metrics

import (
	"sync/atomic"
)

type Collector struct {
	ingested  atomic.Uint64
	stored    atomic.Uint64
	streams   atomic.Int32
	queries   atomic.Uint64
	forwarded atomic.Uint64
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) IncIngested() {
	c.ingested.Add(1)
}

func (c *Collector) IncStored() {
	c.stored.Add(1)
}

func (c *Collector) SetStreams(n int32) {
	c.streams.Store(n)
}

func (c *Collector) IncQueries() {
	c.queries.Add(1)
}

func (c *Collector) IncForwarded() {
	c.forwarded.Add(1)
}

func (c *Collector) Snapshot() Snapshot {
	return Snapshot{
		Ingested:  c.ingested.Load(),
		Stored:    c.stored.Load(),
		Streams:   c.streams.Load(),
		Queries:   c.queries.Load(),
		Forwarded: c.forwarded.Load(),
	}
}

type Snapshot struct {
	Ingested  uint64 `json:"ingested"`
	Stored    uint64 `json:"stored"`
	Streams   int32  `json:"streams"`
	Queries   uint64 `json:"queries"`
	Forwarded uint64 `json:"forwarded"`
}
