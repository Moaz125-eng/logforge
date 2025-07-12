package ingest

import (
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	capacity uint64
	tokens   atomic.Uint64
	last     atomic.Int64
}

func NewRateLimiter(perSecond uint64) *RateLimiter {
	r := &RateLimiter{capacity: perSecond}
	r.tokens.Store(perSecond)
	r.last.Store(time.Now().UnixNano())
	return r
}

func (r *RateLimiter) Allow() bool {
	now := time.Now().UnixNano()
	prev := r.last.Load()
	if now-prev > int64(time.Second) {
		r.last.Store(now)
		r.tokens.Store(r.capacity)
	}
	for {
		cur := r.tokens.Load()
		if cur == 0 {
			return false
		}
		if r.tokens.CompareAndSwap(cur, cur-1) {
			return true
		}
	}
}

func (r *RateLimiter) Remaining() uint64 {
	return r.tokens.Load()
}
