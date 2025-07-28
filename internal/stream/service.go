package stream

import (
	"net/http"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Service struct {
	hub    *Hub
	tail   *TailHandler
	filter *FilteredPublisher
}

func NewService() *Service {
	hub := NewHub()
	return &Service{
		hub:    hub,
		tail:   NewTailHandler(hub),
		filter: NewFilteredPublisher(hub, Filter{}),
	}
}

func (s *Service) Publish(entry logentry.Entry) {
	s.filter.Publish(entry)
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.Handle("/stream/tail", s.tail)
	mux.HandleFunc("/stream/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"clients":` + itoa(uint64(s.hub.Count())) + `}`))
	})
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
