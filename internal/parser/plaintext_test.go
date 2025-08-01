package parser

import "testing"

func TestPlaintextParser(t *testing.T) {
	p := NewPlaintextParser()
	entry, err := p.Parse("simple log line")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if entry.Message != "simple log line" {
		t.Fatalf("unexpected message")
	}
}
