package query

import "testing"

func TestSavedStoreRoundTrip(t *testing.T) {
	store := NewSavedStore()
	store.Put(SavedQuery{Name: "errors", Filter: Filter{Level: "error"}})
	list := store.List()
	if len(list) != 1 {
		t.Fatalf("expected one saved query")
	}
	if !store.Delete(list[0].ID) {
		t.Fatalf("delete failed")
	}
}
