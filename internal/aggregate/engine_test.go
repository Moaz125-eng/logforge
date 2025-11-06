package aggregate

import (
	"testing"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestAggregateCounts(t *testing.T) {
	engine := NewEngine()
	engine.Record(logentry.Entry{Service: "api", Level: logentry.LevelError})
	engine.Record(logentry.Entry{Service: "api", Level: logentry.LevelError})
	snap := engine.Snapshot()
	if snap["total"].(int) != 2 {
		t.Fatalf("expected total 2")
	}
}
