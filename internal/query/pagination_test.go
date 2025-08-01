package query

import (
	"testing"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func TestPaginate(t *testing.T) {
	entries := make([]logentry.Entry, 120)
	page := Paginate(entries, 2, 50)
	if page.Page != 2 {
		t.Fatalf("expected page 2")
	}
	if len(page.Items) != 50 {
		t.Fatalf("expected 50 items")
	}
	if page.Total != 120 {
		t.Fatalf("expected total 120")
	}
}
