package stream

import (
	"strings"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Filter struct {
	Service string
	Level   string
	Keyword string
}

func (f Filter) Allow(entry logentry.Entry) bool {
	if f.Service != "" && entry.Service != f.Service {
		return false
	}
	if f.Level != "" && string(entry.Level) != strings.ToLower(f.Level) {
		return false
	}
	if f.Keyword != "" && !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(f.Keyword)) {
		return false
	}
	return true
}

type FilteredPublisher struct {
	filter Filter
	hub    *Hub
}

func NewFilteredPublisher(hub *Hub, filter Filter) *FilteredPublisher {
	return &FilteredPublisher{filter: filter, hub: hub}
}

func (p *FilteredPublisher) Publish(entry logentry.Entry) {
	if p.filter.Allow(entry) {
		p.hub.Publish(entry)
	}
}
