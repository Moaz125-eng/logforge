package server

import (
	"net/http"

	"github.com/Moaz125-eng/logforge/internal/config"
	"github.com/Moaz125-eng/logforge/internal/ingest"
	"github.com/Moaz125-eng/logforge/internal/index"
	"github.com/Moaz125-eng/logforge/internal/parser"
)

func NewMux(cfg config.Config, ingestSvc *ingest.Service, parserSvc *parser.Service, indexSvc *index.Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler(cfg.NodeID))
	ingestSvc.Register(mux)
	parserSvc.RegisterRoutes(mux)
	indexSvc.Register(mux)
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
