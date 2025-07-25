package storage

import (
	"os"
	"path/filepath"
	"time"
)

type RetentionPolicy struct {
	days int
	dir  string
}

func NewRetentionPolicy(dir string, days int) *RetentionPolicy {
	return &RetentionPolicy{dir: dir, days: days}
}

func (p *RetentionPolicy) Sweep() (int, error) {
	cutoff := time.Now().UTC().Add(-time.Duration(p.days) * 24 * time.Hour)
	removed := 0
	err := filepath.Walk(p.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil {
				return err
			}
			removed++
		}
		return nil
	})
	return removed, err
}

type RetentionWorker struct {
	policy   *RetentionPolicy
	interval time.Duration
	stop     chan struct{}
}

func NewRetentionWorker(policy *RetentionPolicy, interval time.Duration) *RetentionWorker {
	return &RetentionWorker{policy: policy, interval: interval, stop: make(chan struct{})}
}

func (w *RetentionWorker) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-w.stop:
				return
			case <-ticker.C:
				_, _ = w.policy.Sweep()
			}
		}
	}()
}

func (w *RetentionWorker) Stop() {
	close(w.stop)
}
