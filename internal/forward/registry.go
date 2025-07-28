package forward

import (
	"sync"
	"time"
)

type Node struct {
	ID       string    `json:"id"`
	Address  string    `json:"address"`
	LastSeen time.Time `json:"last_seen"`
	Active   bool      `json:"active"`
}

type Registry struct {
	mu    sync.RWMutex
	nodes map[string]Node
}

func NewRegistry() *Registry {
	return &Registry{nodes: make(map[string]Node)}
}

func (r *Registry) Register(node Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	node.LastSeen = time.Now().UTC()
	node.Active = true
	r.nodes[node.ID] = node
}

func (r *Registry) MarkSeen(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if node, ok := r.nodes[id]; ok {
		node.LastSeen = time.Now().UTC()
		node.Active = true
		r.nodes[id] = node
	}
}

func (r *Registry) List() []Node {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Node, 0, len(r.nodes))
	for _, node := range r.nodes {
		out = append(out, node)
	}
	return out
}

func (r *Registry) Prune(maxAge time.Duration) int {
	now := time.Now().UTC()
	removed := 0
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, node := range r.nodes {
		if now.Sub(node.LastSeen) > maxAge {
			delete(r.nodes, id)
			removed++
		}
	}
	return removed
}
