package query

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type SavedQuery struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Filter    Filter    `json:"filter"`
	CreatedAt time.Time `json:"created_at"`
}

type SavedStore struct {
	mu      sync.RWMutex
	queries map[string]SavedQuery
}

func NewSavedStore() *SavedStore {
	return &SavedStore{queries: make(map[string]SavedQuery)}
}

func (s *SavedStore) Put(q SavedQuery) {
	if q.ID == "" {
		q.ID = time.Now().UTC().Format("20060102150405")
	}
	if q.CreatedAt.IsZero() {
		q.CreatedAt = time.Now().UTC()
	}
	s.mu.Lock()
	s.queries[q.ID] = q
	s.mu.Unlock()
}

func (s *SavedStore) Get(id string) (SavedQuery, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.queries[id]
	return q, ok
}

func (s *SavedStore) List() []SavedQuery {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]SavedQuery, 0, len(s.queries))
	for _, q := range s.queries {
		out = append(out, q)
	}
	return out
}

func (s *SavedStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.queries[id]; !ok {
		return false
	}
	delete(s.queries, id)
	return true
}

type SavedService struct {
	store *SavedStore
}

func NewSavedService() *SavedService {
	return &SavedService{store: NewSavedStore()}
}

func (s *SavedService) Register(mux *http.ServeMux) {
	mux.HandleFunc("/query/saved", s.handleSaved)
	mux.HandleFunc("/query/saved/", s.handleOne)
}

func (s *SavedService) handleSaved(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		_ = json.NewEncoder(w).Encode(s.store.List())
	case http.MethodPost:
		var q SavedQuery
		if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		s.store.Put(q)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(q)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *SavedService) handleOne(w http.ResponseWriter, r *http.Request) {
	id := stringsTrimPrefix(r.URL.Path, "/query/saved/")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	if r.Method == http.MethodDelete {
		if !s.store.Delete(id) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	q, ok := s.store.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	_ = json.NewEncoder(w).Encode(q)
}

func stringsTrimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
