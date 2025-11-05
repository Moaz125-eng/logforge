package auth

import "testing"

func TestKeyStoreValidate(t *testing.T) {
	store := NewKeyStore()
	if !store.Validate("logforge-dev-key", "ingest") {
		t.Fatalf("expected bootstrap key valid")
	}
	key := store.Create("ci", []string{"query"})
	if store.Validate(key.Token, "query") != true {
		t.Fatalf("expected new key valid for query")
	}
	store.Revoke(key.ID)
	if store.Validate(key.Token, "query") {
		t.Fatalf("expected revoked key invalid")
	}
}
