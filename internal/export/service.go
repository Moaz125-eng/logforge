package export

import (
	"net/http"
	"strings"

	"github.com/Moaz125-eng/logforge/internal/index"
	"github.com/Moaz125-eng/logforge/internal/query"
)

type Service struct {
	store *index.Store
}

func NewService(store *index.Store) *Service {
	return &Service{store: store}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/export", s.handleExport)
}

func (s *Service) handleExport(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("q")
	entries := s.store.Keyword(term)
	filter := query.Filter{
		Level:   r.URL.Query().Get("level"),
		Service: r.URL.Query().Get("service"),
		Keyword: term,
	}
	entries = query.ApplyFilter(entries, filter)
	format := FormatJSON
	if strings.EqualFold(r.URL.Query().Get("format"), "csv") {
		format = FormatCSV
		w.Header().Set("Content-Type", "text/csv")
	} else {
		w.Header().Set("Content-Type", "application/x-ndjson")
	}
	if err := WriteEntries(w, format, entries); err != nil {
		http.Error(w, "export failed", http.StatusInternalServerError)
	}
}
