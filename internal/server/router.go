package server

import (
	"net/http"

	"github.com/Moaz125-eng/logforge/internal/config"
)

func NewMux(cfg config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler(cfg.NodeID))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"service":"logforge","version":"0.1.0"}`))
	})
	return mux
}
