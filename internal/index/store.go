package index

import (
	"sync"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Store struct {
	inverted  *InvertedIndex
	timestamp *TimestampIndex
	mu        sync.RWMutex
	count     int
}

func NewStore() *Store {
	return &Store{
		inverted:  NewInvertedIndex(),
		timestamp: NewTimestampIndex(),
	}
}

func (s *Store) Add(entry logentry.Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.inverted.Index(entry)
	s.timestamp.Insert(entry)
	s.count++
}

func (s *Store) Keyword(term string) []logentry.Entry {
	return s.inverted.Search(term)
}

func (s *Store) TimeRange(from, to time.Time) []logentry.Entry {
	return s.timestamp.Range(from, to)
}

func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.count
}
