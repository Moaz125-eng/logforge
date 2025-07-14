package parser

import (
	"strings"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type PlaintextParser struct{}

func NewPlaintextParser() *PlaintextParser {
	return &PlaintextParser{}
}

func (p *PlaintextParser) Name() string {
	return "plaintext"
}

func (p *PlaintextParser) Parse(raw string) (logentry.Entry, error) {
	line := strings.TrimSpace(raw)
	entry := logentry.New("plaintext", line, logentry.LevelInfo)
	entry.Raw = raw
	entry.Timestamp = time.Now().UTC()
	parts := strings.SplitN(line, " ", 4)
	if len(parts) >= 4 {
		if ts, err := time.Parse("2006-01-02T15:04:05Z07:00", parts[0]); err == nil {
			entry.Timestamp = ts
			entry.Level = logentry.Level(strings.ToLower(parts[1]))
			entry.Service = parts[2]
			entry.Message = parts[3]
		}
	}
	return entry, nil
}
