package index

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Service struct {
	store *Store
}

func NewService() *Service {
	return &Service{store: NewStore()}
}

func (s *Service) Index(entry logentry.Entry) {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	s.store.Add(entry)
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/index/search", s.searchHandler)
	mux.HandleFunc("/index/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int{"documents": s.store.Count()})
	})
}

func (s *Service) searchHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("q")
	results := s.store.Keyword(term)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"count": len(results), "entries": results})
}

func (s *Service) Store() *Store {
	return s.store
}
