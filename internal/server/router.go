package server

import (
	"net/http"

	"github.com/Moaz125-eng/logforge/internal/aggregate"
	"github.com/Moaz125-eng/logforge/internal/alert"
	"github.com/Moaz125-eng/logforge/internal/auth"
	"github.com/Moaz125-eng/logforge/internal/config"
	"github.com/Moaz125-eng/logforge/internal/export"
	"github.com/Moaz125-eng/logforge/internal/forward"
	"github.com/Moaz125-eng/logforge/internal/metrics"
	"github.com/Moaz125-eng/logforge/internal/ingest"
	"github.com/Moaz125-eng/logforge/internal/index"
	"github.com/Moaz125-eng/logforge/internal/parser"
	"github.com/Moaz125-eng/logforge/internal/query"
	"github.com/Moaz125-eng/logforge/internal/storage"
	"github.com/Moaz125-eng/logforge/internal/stream"
)

func NewMux(cfg config.Config, ingestSvc *ingest.Service, parserSvc *parser.Service, indexSvc *index.Service, queryEngine *query.Engine, savedSvc *query.SavedService, storageSvc *storage.Service, streamSvc *stream.Service, forwardSvc *forward.Service, metricsSvc *metrics.Service, alertSvc *alert.Service, authSvc *auth.Service, aggregateSvc *aggregate.Service, exportSvc *export.Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler(cfg.NodeID))
	guard := func(h http.Handler) http.Handler {
		return authSvc.Middleware().Require("ingest", h)
	}
	ingestSvc.Register(mux, guard)
	authSvc.Register(mux)
	parserSvc.RegisterRoutes(mux)
	indexSvc.Register(mux)
	queryEngine.Register(mux)
	savedSvc.Register(mux)
	storageSvc.Register(mux)
	streamSvc.Register(mux)
	forwardSvc.Register(mux)
	metricsSvc.Register(mux)
	alertSvc.Register(mux)
	aggregateSvc.Register(mux)
	exportSvc.Register(mux)
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
