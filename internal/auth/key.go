package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type APIKey struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	Name      string    `json:"name"`
	Scopes    []string  `json:"scopes"`
	CreatedAt time.Time `json:"created_at"`
	Revoked   bool      `json:"revoked"`
}

type KeyStore struct {
	mu   sync.RWMutex
	keys map[string]APIKey
}

func NewKeyStore() *KeyStore {
	store := &KeyStore{keys: make(map[string]APIKey)}
	store.keys["bootstrap"] = APIKey{
		ID: "bootstrap", Token: "logforge-dev-key", Name: "development",
		Scopes: []string{"ingest", "query"}, CreatedAt: time.Now().UTC(),
	}
	return store
}

func (s *KeyStore) Create(name string, scopes []string) APIKey {
	token := randomToken()
	key := APIKey{
		ID: randomID(), Token: token, Name: name, Scopes: scopes,
		CreatedAt: time.Now().UTC(),
	}
	s.mu.Lock()
	s.keys[key.ID] = key
	s.mu.Unlock()
	return key
}

func (s *KeyStore) Revoke(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	key, ok := s.keys[id]
	if !ok {
		return false
	}
	key.Revoked = true
	s.keys[id] = key
	return true
}

func (s *KeyStore) Validate(token string, scope string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, key := range s.keys {
		if key.Revoked || key.Token != token {
			continue
		}
		for _, sc := range key.Scopes {
			if sc == scope || sc == "admin" {
				return true
			}
		}
	}
	return false
}

func (s *KeyStore) List() []APIKey {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]APIKey, 0, len(s.keys))
	for _, key := range s.keys {
		copy := key
		copy.Token = mask(copy.Token)
		out = append(out, copy)
	}
	return out
}

func mask(token string) string {
	if len(token) <= 6 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-2:]
}

func randomToken() string {
	buf := make([]byte, 24)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

func randomID() string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
