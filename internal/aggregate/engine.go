package aggregate

import (
	"sync"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Bucket struct {
	Service string `json:"service"`
	Level   string `json:"level"`
	Count   int    `json:"count"`
}

type Engine struct {
	mu       sync.RWMutex
	counts   map[string]int
	services map[string]int
	levels   map[string]int
	window   time.Time
}

func NewEngine() *Engine {
	return &Engine{
		counts:   make(map[string]int),
		services: make(map[string]int),
		levels:   make(map[string]int),
		window:   time.Now().UTC(),
	}
}

func (e *Engine) Record(entry logentry.Entry) {
	e.mu.Lock()
	defer e.mu.Unlock()
	key := entry.Service + "|" + string(entry.Level)
	e.counts[key]++
	e.services[entry.Service]++
	e.levels[string(entry.Level)]++
}

func (e *Engine) Snapshot() map[string]any {
	e.mu.RLock()
	defer e.mu.RUnlock()
	byPair := make([]Bucket, 0, len(e.counts))
	for key, count := range e.counts {
		service, level := splitKey(key)
		byPair = append(byPair, Bucket{Service: service, Level: level, Count: count})
	}
	return map[string]any{
		"window_start": e.window,
		"by_pair":      byPair,
		"by_service":   e.services,
		"by_level":     e.levels,
		"total":        sumMap(e.counts),
	}
}

func (e *Engine) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.counts = make(map[string]int)
	e.services = make(map[string]int)
	e.levels = make(map[string]int)
	e.window = time.Now().UTC()
}

func splitKey(key string) (string, string) {
	for i := 0; i < len(key); i++ {
		if key[i] == '|' {
			return key[:i], key[i+1:]
		}
	}
	return key, ""
}

func sumMap(m map[string]int) int {
	total := 0
	for _, v := range m {
		total += v
	}
	return total
}
