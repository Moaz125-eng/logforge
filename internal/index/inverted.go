package index

import (
	"strings"
	"sync"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type InvertedIndex struct {
	mu     sync.RWMutex
	terms  map[string]map[string]struct{}
	ids    map[string]logentry.Entry
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		terms: make(map[string]map[string]struct{}),
		ids:   make(map[string]logentry.Entry),
	}
}

func (idx *InvertedIndex) Index(entry logentry.Entry) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	id := entry.ID
	if id == "" {
		id = entry.Message
	}
	idx.ids[id] = entry
	for _, term := range tokenize(entry.Message) {
		bucket, ok := idx.terms[term]
		if !ok {
			bucket = make(map[string]struct{})
			idx.terms[term] = bucket
		}
		bucket[id] = struct{}{}
	}
	for k, v := range entry.Fields {
		for _, term := range tokenize(k + ":" + v) {
			bucket, ok := idx.terms[term]
			if !ok {
				bucket = make(map[string]struct{})
				idx.terms[term] = bucket
			}
			bucket[id] = struct{}{}
		}
	}
}

func (idx *InvertedIndex) Search(term string) []logentry.Entry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	bucket := idx.terms[strings.ToLower(term)]
	if bucket == nil {
		return nil
	}
	out := make([]logentry.Entry, 0, len(bucket))
	for id := range bucket {
		if entry, ok := idx.ids[id]; ok {
			out = append(out, entry)
		}
	}
	return out
}

func tokenize(text string) []string {
	parts := strings.Fields(strings.ToLower(text))
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if len(p) > 1 {
			out = append(out, p)
		}
	}
	return out
}
