package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestGzipRoundTrip(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "logforge-test")
	store, err := NewGzipStore(dir)
	if err != nil {
		t.Fatalf("store init failed: %v", err)
	}
	entries := []logentry.Entry{
		logentry.New("a", "one", logentry.LevelInfo),
		logentry.New("b", "two", logentry.LevelWarn),
	}
	path, err := store.Persist(entries)
	if err != nil {
		t.Fatalf("persist failed: %v", err)
	}
	loaded, err := store.Read(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries")
	}
}
