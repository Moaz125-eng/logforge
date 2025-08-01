package bench

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Result struct {
	Requests uint64
	Errors   uint64
	Duration time.Duration
	LatencyP50 time.Duration
	LatencyP95 time.Duration
}

type Runner struct {
	target     string
	workers    int
	total      int
	latencies  []time.Duration
	mu         sync.Mutex
	errors     atomic.Uint64
	done       atomic.Uint64
}

func NewRunner(target string, workers, total int) *Runner {
	return &Runner{target: target, workers: workers, total: total, latencies: make([]time.Duration, 0, total)}
}

func (r *Runner) Run(ctx context.Context) Result {
	start := time.Now()
	wg := sync.WaitGroup{}
	jobs := make(chan int, r.workers*2)
	for i := 0; i < r.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				if ctx.Err() != nil {
					return
				}
				latency, err := r.sendOnce(ctx)
				r.mu.Lock()
				r.latencies = append(r.latencies, latency)
				r.mu.Unlock()
				if err != nil {
					r.errors.Add(1)
				}
				r.done.Add(1)
			}
		}()
	}
	for i := 0; i < r.total; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	duration := time.Since(start)
	return Result{
		Requests: r.done.Load(),
		Errors:   r.errors.Load(),
		Duration: duration,
		LatencyP50: percentile(r.latencies, 0.50),
		LatencyP95: percentile(r.latencies, 0.95),
	}
}

func (r *Runner) sendOnce(ctx context.Context) (time.Duration, error) {
	entry := logentry.New("bench", fmt.Sprintf("event-%d", time.Now().UnixNano()), logentry.LevelInfo)
	data, err := json.Marshal(entry)
	if err != nil {
		return 0, err
	}
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.target+"/ingest", bytes.NewReader(data))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return time.Since(start), nil
}

func percentile(samples []time.Duration, p float64) time.Duration {
	if len(samples) == 0 {
		return 0
	}
	idx := int(float64(len(samples)-1) * p)
	if idx < 0 {
		idx = 0
	}
	return samples[idx]
}
