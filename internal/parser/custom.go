package parser

import (
	"regexp"
	"strings"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type CustomParser struct {
	name    string
	pattern *regexp.Regexp
	groups  map[string]int
}

func NewCustomParser(name, pattern string, groups map[string]int) (*CustomParser, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &CustomParser{name: name, pattern: re, groups: groups}, nil
}

func (p *CustomParser) Name() string {
	return p.name
}

func (p *CustomParser) Parse(raw string) (logentry.Entry, error) {
	matches := p.pattern.FindStringSubmatch(raw)
	entry := logentry.New(p.name, raw, logentry.LevelInfo)
	entry.Raw = raw
	entry.Timestamp = time.Now().UTC()
	if matches == nil {
		return entry, nil
	}
	if idx, ok := p.groups["service"]; ok && idx < len(matches) {
		entry.Service = matches[idx]
	}
	if idx, ok := p.groups["level"]; ok && idx < len(matches) {
		entry.Level = logentry.Level(strings.ToLower(matches[idx]))
	}
	if idx, ok := p.groups["message"]; ok && idx < len(matches) {
		entry.Message = matches[idx]
	}
	return entry, nil
}

type CustomBuilder struct {
	registry *Registry
}

func NewCustomBuilder(registry *Registry) *CustomBuilder {
	return &CustomBuilder{registry: registry}
}

func (b *CustomBuilder) Add(name, pattern string, groups map[string]int) error {
	p, err := NewCustomParser(name, pattern, groups)
	if err != nil {
		return err
	}
	b.registry.Register(p)
	return nil
}
