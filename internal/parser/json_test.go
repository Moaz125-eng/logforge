package parser

import "testing"

func TestJSONParserParsesPayload(t *testing.T) {
	p := NewJSONParser()
	raw := `{"timestamp":"2025-07-12T10:00:00Z","level":"error","service":"api","message":"failed","fields":{"region":"us"}}`
	entry, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if entry.Service != "api" {
		t.Fatalf("expected api service")
	}
	if entry.Message != "failed" {
		t.Fatalf("expected failed message")
	}
}
