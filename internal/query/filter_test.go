package query

import (
	"testing"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestFilterMatch(t *testing.T) {
	entry := logentry.Entry{
		Level:     logentry.LevelError,
		Service:   "billing",
		Message:   "charge failed",
		Timestamp: time.Now().UTC(),
	}
	filter := Filter{Level: "error", Service: "billing", Keyword: "charge"}
	if !filter.Match(entry) {
		t.Fatalf("expected match")
	}
}

func TestApplyFilter(t *testing.T) {
	entries := []logentry.Entry{
		{Level: logentry.LevelInfo, Service: "a", Message: "ok"},
		{Level: logentry.LevelError, Service: "b", Message: "bad"},
	}
	filtered := ApplyFilter(entries, Filter{Level: "error"})
	if len(filtered) != 1 {
		t.Fatalf("expected one filtered entry")
	}
}
