package alert

import (
	"sync"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Engine struct {
	mu     sync.RWMutex
	rules  map[string]Rule
	events []Event
	limit  int
}

func NewEngine(limit int) *Engine {
	return &Engine{rules: make(map[string]Rule), events: make([]Event, 0, limit), limit: limit}
}

func (e *Engine) Upsert(rule Rule) {
	if rule.Window == 0 {
		rule.Window = time.Minute
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules[rule.ID] = rule
}

func (e *Engine) Delete(id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.rules, id)
}

func (e *Engine) List() []Rule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]Rule, 0, len(e.rules))
	for _, rule := range e.rules {
		out = append(out, rule)
	}
	return out
}

func (e *Engine) Evaluate(entry logentry.Entry) []Event {
	now := time.Now().UTC()
	e.mu.Lock()
	defer e.mu.Unlock()
	fired := make([]Event, 0)
	for id, rule := range e.rules {
		if !rule.Match(entry) {
			continue
		}
		copy := rule
		if copy.Bump(now) {
			ev := NewEvent(copy, entry)
			fired = append(fired, ev)
			e.rules[id] = copy
			e.events = append(e.events, ev)
			if len(e.events) > e.limit {
				e.events = e.events[len(e.events)-e.limit:]
			}
		} else {
			e.rules[id] = copy
		}
	}
	return fired
}

func (e *Engine) Recent() []Event {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]Event, len(e.events))
	copy(out, e.events)
	return out
}
