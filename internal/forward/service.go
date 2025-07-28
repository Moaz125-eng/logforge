package forward

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Moaz125-eng/logforge/internal/config"
)

type Service struct {
	registry *Registry
	agent    *Agent
	failover *Failover
}

func NewService(cfg config.Config) *Service {
	registry := NewRegistry()
	registry.Register(Node{ID: cfg.NodeID, Address: cfg.HTTPAddr, Active: true})
	var primary string
	if len(cfg.ForwardPeers) > 0 {
		primary = cfg.ForwardPeers[0]
	}
	backups := cfg.ForwardPeers
	agent := NewAgent(cfg.ForwardPeers)
	failover := NewFailover(primary, backups)
	return &Service{registry: registry, agent: agent, failover: failover}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
		s.registry.Prune(2 * time.Minute)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(s.registry.List())
	})
	mux.HandleFunc("/nodes/register", func(w http.ResponseWriter, r *http.Request) {
		var node Node
		if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		s.registry.Register(node)
		w.WriteHeader(http.StatusCreated)
	})
}

func (s *Service) Registry() *Registry {
	return s.registry
}

func (s *Service) Agent() *Agent {
	return s.agent
}
