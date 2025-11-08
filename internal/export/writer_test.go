package export

import (
	"bytes"
	"testing"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestExportJSON(t *testing.T) {
	var buf bytes.Buffer
	entry := logentry.Entry{
		Timestamp: time.Now().UTC(),
		Level:     logentry.LevelInfo,
		Service:   "api",
		Message:   "ok",
	}
	if err := WriteEntries(&buf, FormatJSON, []logentry.Entry{entry}); err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("api")) {
		t.Fatalf("expected service in output")
	}
}
