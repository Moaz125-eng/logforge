package metrics

import (
	"encoding/json"
	"net/http"
)

type Dashboard struct {
	collector  *Collector
	prometheus *Prometheus
}

func NewDashboard(collector *Collector, prometheus *Prometheus) *Dashboard {
	return &Dashboard{collector: collector, prometheus: prometheus}
}

func (d *Dashboard) Register(mux *http.ServeMux) {
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		snap := d.collector.Snapshot()
		d.prometheus.Observe(snap)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(snap)
	})
	mux.Handle("/metrics/prometheus", d.prometheus.Handler())
}
