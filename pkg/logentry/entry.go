package logentry

import "time"

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

type Entry struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	Service   string            `json:"service"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	Raw       string            `json:"raw,omitempty"`
}

func New(service, message string, level Level) Entry {
	return Entry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Service:   service,
		Message:   message,
		Fields:    make(map[string]string),
	}
}
