package alert

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Service struct {
	engine   *Engine
	notifier *Notifier
}

func NewService(webhook string) *Service {
	engine := NewEngine(500)
	engine.Upsert(Rule{
		ID: "default-error", Name: "error spike", Level: "error",
		Severity: SeverityWarn, Threshold: 5, Window: time.Minute,
	})
	return &Service{engine: engine, notifier: NewNotifier(webhook)}
}

func (s *Service) Watch(entry logentry.Entry) {
	events := s.engine.Evaluate(entry)
	s.notifier.Notify(events)
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/alerts/rules", s.handleRules)
	mux.HandleFunc("/alerts/events", s.handleEvents)
}

func (s *Service) handleRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		_ = json.NewEncoder(w).Encode(s.engine.List())
	case http.MethodPost:
		var rule Rule
		if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
			http.Error(w, "invalid rule", http.StatusBadRequest)
			return
		}
		if rule.ID == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		s.engine.Upsert(rule)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(rule)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		s.engine.Delete(id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Service) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"events": s.engine.Recent(),
		"sent":   s.notifier.Sent(),
	})
}
