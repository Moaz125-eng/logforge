package ingest

import (
	"io"
	"net/http"
	"sync/atomic"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Sink func(logentry.Entry) error

type Handler struct {
	sink      Sink
	accepted  atomic.Uint64
	rejected  atomic.Uint64
	maxBody   int64
}

func NewHandler(sink Sink) *Handler {
	return &Handler{sink: sink, maxBody: 1 << 20}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	body, err := io.ReadAll(io.LimitReader(r.Body, h.maxBody))
	if err != nil {
		h.rejected.Add(1)
		http.Error(w, "read failed", http.StatusBadRequest)
		return
	}
	entry := logentry.Entry{
		Service: r.Header.Get("X-Log-Service"),
		Message: string(body),
		Raw:     string(body),
	}
	if entry.Service == "" {
		entry.Service = "http"
	}
	if err := h.sink(entry); err != nil {
		h.rejected.Add(1)
		http.Error(w, "ingest failed", http.StatusInternalServerError)
		return
	}
	h.accepted.Add(1)
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) Stats() (accepted, rejected uint64) {
	return h.accepted.Load(), h.rejected.Load()
}
