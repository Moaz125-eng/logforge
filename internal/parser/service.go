package parser

import (
	"net/http"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Service struct {
	registry *Registry
	parsed   uint64
	failed   uint64
}

func NewService() *Service {
	return &Service{registry: NewRegistry()}
}

func (s *Service) Parse(format, raw string) (logentry.Entry, error) {
	entry, err := s.registry.Parse(format, raw)
	if err != nil {
		s.failed++
		return logentry.Entry{}, err
	}
	s.parsed++
	return entry, nil
}

func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/parsers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		names := s.registry.Names()
		w.Write([]byte(`{"parsers":["` + join(names, `","`) + `"]}`))
	})
}

func join(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += sep + parts[i]
	}
	return out
}
