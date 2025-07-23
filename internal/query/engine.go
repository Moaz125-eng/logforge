package query

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Moaz125-eng/logforge/internal/index"
)

type Engine struct {
	index *index.Store
}

func NewEngine(store *index.Store) *Engine {
	return &Engine{index: store}
}

func (e *Engine) Register(mux *http.ServeMux) {
	mux.HandleFunc("/query", e.handleQuery)
}

func (e *Engine) handleQuery(w http.ResponseWriter, r *http.Request) {
	filter := Filter{
		Level:   r.URL.Query().Get("level"),
		Service: r.URL.Query().Get("service"),
		Keyword: r.URL.Query().Get("q"),
	}
	if raw := r.URL.Query().Get("from"); raw != "" {
		if ts, err := time.Parse(time.RFC3339, raw); err == nil {
			filter.From = ts
		}
	}
	if raw := r.URL.Query().Get("to"); raw != "" {
		if ts, err := time.Parse(time.RFC3339, raw); err == nil {
			filter.To = ts
		}
	}
	entries := e.index.Keyword(filter.Keyword)
	if !filter.From.IsZero() || !filter.To.IsZero() {
		to := filter.To
		if to.IsZero() {
			to = time.Now().UTC()
		}
		from := filter.From
		if from.IsZero() {
			from = time.Time{}
		}
		entries = append(entries, e.index.TimeRange(from, to)...)
	}
	filtered := ApplyFilter(entries, filter)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	result := Paginate(filtered, page, size)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
