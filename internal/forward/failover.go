package forward

import (
	"sync"
)

type Failover struct {
	mu       sync.RWMutex
	primary  string
	backups  []string
	active   string
	failures map[string]int
}

func NewFailover(primary string, backups []string) *Failover {
	return &Failover{
		primary:  primary,
		backups:  backups,
		active:   primary,
		failures: make(map[string]int),
	}
}

func (f *Failover) Active() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.active
}

func (f *Failover) ReportFailure(target string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.failures[target]++
	if f.active == target && f.failures[target] >= 2 {
		for _, backup := range f.backups {
			if f.failures[backup] < 2 {
				f.active = backup
				return
			}
		}
	}
}

func (f *Failover) ReportSuccess(target string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.failures[target] = 0
	if f.active != f.primary {
		f.active = f.primary
	}
}
