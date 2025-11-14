package config

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

type Reloader struct {
	mu     sync.RWMutex
	active Config
}

func NewReloader(initial Config) *Reloader {
	return &Reloader{active: initial}
}

func (r *Reloader) Current() Config {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.active
}

func (r *Reloader) Reload() Config {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.active = Load()
	return r.active
}

func (r *Reloader) ApplyPatch(patch map[string]string) Config {
	r.mu.Lock()
	defer r.mu.Unlock()
	for key, value := range patch {
		_ = os.Setenv(key, value)
	}
	r.active = Load()
	return r.active
}

func (r *Reloader) Register(mux *http.ServeMux) {
	mux.HandleFunc("/config/reload", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch req.Method {
		case http.MethodPost:
			var patch map[string]string
			if req.Body != nil && req.ContentLength > 0 {
				_ = json.NewDecoder(req.Body).Decode(&patch)
				_ = json.NewEncoder(w).Encode(r.ApplyPatch(patch))
				return
			}
			_ = json.NewEncoder(w).Encode(r.Reload())
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(r.Current())
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
