package index

import (
	"testing"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestInvertedIndexSearch(t *testing.T) {
	idx := NewInvertedIndex()
	entry := logentry.New("api", "user login success", logentry.LevelInfo)
	entry.ID = "1"
	idx.Index(entry)
	results := idx.Search("login")
	if len(results) != 1 {
		t.Fatalf("expected one result got %d", len(results))
	}
}
