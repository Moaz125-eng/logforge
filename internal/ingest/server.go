package ingest

import (
	"context"
	"net/http"

	"github.com/Moaz125-eng/logforge/internal/config"
	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Service struct {
	cfg    config.Config
	http   *Handler
	tcp    *TCPServer
	dlq    *DLQ
}

func NewService(cfg config.Config, sink Sink) *Service {
	httpHandler := NewHandler(sink)
	tcpServer := NewTCPServer(cfg.TCPAddr, sink)
	return &Service{cfg: cfg, http: httpHandler, tcp: tcpServer, dlq: NewDLQ(1000, 5)}
}

func (s *Service) HTTPHandler() http.Handler {
	return s.http
}

func (s *Service) Register(mux *http.ServeMux, guard func(http.Handler) http.Handler) {
	handler := s.http
	if guard != nil {
		handler = guard(handler)
	}
	mux.Handle("/ingest", handler)
	s.dlq.Register(mux, s.replaySink())
	mux.HandleFunc("/ingest/stats", func(w http.ResponseWriter, r *http.Request) {
		a, rej := s.http.Stats()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"accepted":` + itoa(a) + `,"rejected":` + itoa(rej) + `,"tcp_lines":` + itoa(s.tcp.LinesIngested()) + `}`))
	})
}

func (s *Service) Start(ctx context.Context) error {
	return s.tcp.Start(ctx)
}

func (s *Service) Wait() {
	s.tcp.Wait()
}

func (s *Service) DLQ() *DLQ {
	return s.dlq
}

func (s *Service) replaySink() func(logentry.Entry) error {
	return func(entry logentry.Entry) error {
		return s.http.sink(entry)
	}
}

func noopSink(e logentry.Entry) error {
	return nil
}

func itoa(v uint64) string {
	if v == 0 {
		return "0"
	}
	buf := make([]byte, 20)
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
