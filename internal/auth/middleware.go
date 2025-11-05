package auth

import (
	"net/http"
	"strings"
)

type Middleware struct {
	store *KeyStore
}

func NewMiddleware(store *KeyStore) *Middleware {
	return &Middleware{store: store}
}

func (m *Middleware) Require(scope string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" || !m.store.Validate(token, scope) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Optional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token != "" && !m.store.Validate(token, "ingest") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func extractToken(r *http.Request) string {
	if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return r.Header.Get("X-API-Key")
}

func WrapHandler(store *KeyStore, scope string, handler http.HandlerFunc) http.Handler {
	return NewMiddleware(store).Require(scope, handler)
}
