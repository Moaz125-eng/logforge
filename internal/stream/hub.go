package stream

import (
	"sync"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[chan logentry.Entry]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[chan logentry.Entry]struct{})}
}

func (h *Hub) Subscribe(buffer int) chan logentry.Entry {
	ch := make(chan logentry.Entry, buffer)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(ch chan logentry.Entry) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
	close(ch)
}

func (h *Hub) Publish(entry logentry.Entry) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.clients {
		select {
		case ch <- entry:
		default:
		}
	}
}

func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
