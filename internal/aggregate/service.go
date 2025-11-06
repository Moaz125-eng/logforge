package aggregate

import (
	"encoding/json"
	"net/http"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Service struct {
	engine *Engine
}

func NewService() *Service {
	return &Service{engine: NewEngine()}
}

func (s *Service) Record(entry logentry.Entry) {
	s.engine.Record(entry)
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/aggregate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodDelete {
			s.engine.Reset()
			w.WriteHeader(http.StatusNoContent)
			return
		}
		_ = json.NewEncoder(w).Encode(s.engine.Snapshot())
	})
}
