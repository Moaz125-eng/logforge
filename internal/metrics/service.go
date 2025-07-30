package metrics

import "net/http"

type Service struct {
	collector  *Collector
	prometheus *Prometheus
	dashboard  *Dashboard
}

func NewService() *Service {
	collector := NewCollector()
	prometheus := NewPrometheus()
	dashboard := NewDashboard(collector, prometheus)
	return &Service{collector: collector, prometheus: prometheus, dashboard: dashboard}
}

func (s *Service) Register(mux *http.ServeMux) {
	s.dashboard.Register(mux)
}

func (s *Service) Collector() *Collector {
	return s.collector
}
