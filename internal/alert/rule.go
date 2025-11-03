package alert

import (
	"strings"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Severity string

const (
	SeverityInfo  Severity = "info"
	SeverityWarn  Severity = "warn"
	SeverityCrit  Severity = "crit"
)

type Rule struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Service   string   `json:"service"`
	Level     string   `json:"level"`
	Contains  string   `json:"contains"`
	Severity  Severity `json:"severity"`
	Threshold int      `json:"threshold"`
	Window    time.Duration
	hits      int
	windowEnd time.Time
}

func (r *Rule) Match(entry logentry.Entry) bool {
	if r.Service != "" && entry.Service != r.Service {
		return false
	}
	if r.Level != "" && string(entry.Level) != strings.ToLower(r.Level) {
		return false
	}
	if r.Contains != "" && !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(r.Contains)) {
		return false
	}
	return true
}

func (r *Rule) Bump(now time.Time) bool {
	if r.Threshold <= 1 {
		return true
	}
	if now.After(r.windowEnd) {
		r.hits = 0
		r.windowEnd = now.Add(r.Window)
	}
	r.hits++
	return r.hits >= r.Threshold
}

type Event struct {
	RuleID    string    `json:"rule_id"`
	Name      string    `json:"name"`
	Severity  Severity  `json:"severity"`
	Message   string    `json:"message"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

func NewEvent(rule Rule, entry logentry.Entry) Event {
	return Event{
		RuleID:    rule.ID,
		Name:      rule.Name,
		Severity:  rule.Severity,
		Message:   entry.Message,
		Service:   entry.Service,
		Timestamp: time.Now().UTC(),
	}
}
