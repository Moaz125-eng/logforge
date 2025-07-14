package parser

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type JSONParser struct{}

func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

func (p *JSONParser) Name() string {
	return "json"
}

type jsonPayload struct {
	Timestamp string            `json:"timestamp"`
	Level     string            `json:"level"`
	Service   string            `json:"service"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields"`
}

func (p *JSONParser) Parse(raw string) (logentry.Entry, error) {
	var payload jsonPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return logentry.Entry{}, err
	}
	entry := logentry.New(payload.Service, payload.Message, logentry.LevelInfo)
	if payload.Level != "" {
		entry.Level = logentry.Level(strings.ToLower(payload.Level))
	}
	if payload.Timestamp != "" {
		if ts, err := time.Parse(time.RFC3339, payload.Timestamp); err == nil {
			entry.Timestamp = ts
		}
	}
	if len(payload.Fields) > 0 {
		entry.Fields = payload.Fields
	}
	entry.Raw = raw
	return entry, nil
}
