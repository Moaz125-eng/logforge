package auth

import (
	"encoding/json"
	"net/http"
)

type Service struct {
	store *KeyStore
	mw    *Middleware
}

func NewService() *Service {
	store := NewKeyStore()
	return &Service{store: store, mw: NewMiddleware(store)}
}

func (s *Service) Store() *KeyStore {
	return s.store
}

func (s *Service) Middleware() *Middleware {
	return s.mw
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/auth/keys", s.handleKeys)
}

func (s *Service) handleKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		_ = json.NewEncoder(w).Encode(s.store.List())
	case http.MethodPost:
		var body struct {
			Name   string   `json:"name"`
			Scopes []string `json:"scopes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			body.Name = "unnamed"
		}
		if len(body.Scopes) == 0 {
			body.Scopes = []string{"ingest"}
		}
		key := s.store.Create(body.Name, body.Scopes)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(key)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !s.store.Revoke(id) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
