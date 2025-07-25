package storage

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Moaz125-eng/logforge/internal/config"
	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Service struct {
	manager *Manager
	worker  *RetentionWorker
}

func NewService(cfg config.Config) (*Service, error) {
	manager, err := NewManager(cfg.DataDir)
	if err != nil {
		return nil, err
	}
	policy := NewRetentionPolicy(cfg.DataDir, cfg.RetentionDays)
	worker := NewRetentionWorker(policy, time.Hour)
	worker.Start()
	return &Service{manager: manager, worker: worker}, nil
}

func (s *Service) Persist(entry logentry.Entry) error {
	return s.manager.Store(entry)
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/storage/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int64{"bytes": s.manager.BytesStored()})
	})
}

func (s *Service) Close() error {
	s.worker.Stop()
	return s.manager.Close()
}
