package parser

import (
	"sync"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Parser interface {
	Name() string
	Parse(raw string) (logentry.Entry, error)
}

type Registry struct {
	mu      sync.RWMutex
	parsers map[string]Parser
	order   []string
}

func NewRegistry() *Registry {
	r := &Registry{parsers: make(map[string]Parser)}
	r.Register(NewJSONParser())
	r.Register(NewPlaintextParser())
	return r
}

func (r *Registry) Register(p Parser) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parsers[p.Name()] = p
	r.order = append(r.order, p.Name())
}

func (r *Registry) Parse(format, raw string) (logentry.Entry, error) {
	r.mu.RLock()
	p, ok := r.parsers[format]
	r.mu.RUnlock()
	if !ok {
		p = r.parsers["plaintext"]
	}
	return p.Parse(raw)
}

func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, len(r.order))
	copy(out, r.order)
	return out
}
