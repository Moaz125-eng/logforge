package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	ingested  prometheus.Counter
	stored    prometheus.Counter
	streams   prometheus.Gauge
	queries   prometheus.Counter
	forwarded prometheus.Counter
}

func NewPrometheus() *Prometheus {
	p := &Prometheus{
		ingested: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "logforge_ingested_total",
			Help: "Total ingested log entries",
		}),
		stored: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "logforge_stored_total",
			Help: "Total stored log entries",
		}),
		streams: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "logforge_active_streams",
			Help: "Active websocket stream clients",
		}),
		queries: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "logforge_queries_total",
			Help: "Total query requests",
		}),
		forwarded: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "logforge_forwarded_total",
			Help: "Total forwarded log entries",
		}),
	}
	prometheus.MustRegister(p.ingested, p.stored, p.streams, p.queries, p.forwarded)
	return p
}

func (p *Prometheus) Handler() http.Handler {
	return promhttp.Handler()
}

func (p *Prometheus) Observe(snapshot Snapshot) {
	p.ingested.Add(float64(snapshot.Ingested))
	p.stored.Add(float64(snapshot.Stored))
	p.streams.Set(float64(snapshot.Streams))
	p.queries.Add(float64(snapshot.Queries))
	p.forwarded.Add(float64(snapshot.Forwarded))
}
