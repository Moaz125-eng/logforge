package pipeline

import (
	"context"
	"testing"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestPipelineRunsStages(t *testing.T) {
	pipe := NewPipeline(&EnrichStage{})
	entry := logentry.New("svc", "msg", logentry.LevelInfo)
	out, err := pipe.Run(context.Background(), entry)
	if err != nil {
		t.Fatalf("pipeline failed: %v", err)
	}
	if out.Fields["node"] != "local" {
		t.Fatalf("expected enrich field")
	}
}
