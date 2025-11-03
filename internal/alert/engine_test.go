package alert

import (
	"testing"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestRuleMatchAndThreshold(t *testing.T) {
	rule := Rule{
		ID: "t", Name: "test", Level: "error", Severity: SeverityCrit,
		Threshold: 2, Window: time.Minute,
	}
	entry := logentry.Entry{Level: logentry.LevelError, Service: "api", Message: "boom"}
	if !rule.Match(entry) {
		t.Fatalf("expected match")
	}
	now := time.Now().UTC()
	if !rule.Bump(now) {
		t.Fatalf("expected first bump")
	}
	if !rule.Bump(now) {
		t.Fatalf("expected threshold fire")
	}
}

func TestEngineEvaluate(t *testing.T) {
	engine := NewEngine(10)
	engine.Upsert(Rule{ID: "e", Name: "err", Level: "error", Severity: SeverityWarn, Threshold: 1})
	events := engine.Evaluate(logentry.Entry{Level: logentry.LevelError, Message: "x"})
	if len(events) != 1 {
		t.Fatalf("expected one event")
	}
}
