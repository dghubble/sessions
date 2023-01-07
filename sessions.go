package sessions

import (
	"net/http"
)

const (
	defaultMaxAge = 3600 * 24 * 7 // 1 week
)

// Session represents state values maintained in a sessions Store.
type Session[V any] struct {
	name   string
	values map[string]V
	// convenience methods Save and Destroy use store
	store Store[V]
}

// NewSession returns a new Session.
func NewSession[V any](store Store[V], name string) *Session[V] {
	return &Session[V]{
		name:   name,
		values: make(map[string]V),
		store:  store,
	}
}

// Name returns the name of the session.
func (s *Session[V]) Name() string {
	return s.name
}

// Set sets a key/value pair in the session state.
func (s *Session[V]) Set(key string, value V) {
	s.values[key] = value
}

// Get returns the state value for the given key.
func (s *Session[V]) Get(key string) V {
	return s.values[key]
}

// GetOk returns the state value for the given key and whether they key exists.
func (s *Session[V]) GetOk(key string) (V, bool) {
	value, ok := s.values[key]
	return value, ok
}

// Save adds or updates the session. Identical to calling
// store.Save(w, session).
func (s *Session[V]) Save(w http.ResponseWriter) error {
	return s.store.Save(w, s)
}

// Destroy destroys the session. Identical to calling
// store.Destroy(w, session.name).
func (s *Session[V]) Destroy(w http.ResponseWriter) {
	s.store.Destroy(w, s.name)
}
