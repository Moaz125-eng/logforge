package query

import (
	"strings"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Filter struct {
	Level   string
	Service string
	Keyword string
	From    time.Time
	To      time.Time
}

func (f Filter) Match(entry logentry.Entry) bool {
	if f.Level != "" && string(entry.Level) != strings.ToLower(f.Level) {
		return false
	}
	if f.Service != "" && entry.Service != f.Service {
		return false
	}
	if f.Keyword != "" && !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(f.Keyword)) {
		return false
	}
	if !f.From.IsZero() && entry.Timestamp.Before(f.From) {
		return false
	}
	if !f.To.IsZero() && entry.Timestamp.After(f.To) {
		return false
	}
	return true
}

func ApplyFilter(entries []logentry.Entry, filter Filter) []logentry.Entry {
	out := make([]logentry.Entry, 0, len(entries))
	for _, entry := range entries {
		if filter.Match(entry) {
			out = append(out, entry)
		}
	}
	return out
}
