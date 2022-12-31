package sessions

import (
	"net/http"
)

const (
	defaultMaxAge = 3600 * 24 * 7 // 1 week
)

// Session represents state values maintained in a sessions Store.
type Session struct {
	name   string
	values map[string]any
	// convenience methods Save and Destroy use store
	store Store
}

// NewSession returns a new Session.
func NewSession(store Store, name string) *Session {
	return &Session{
		name:   name,
		values: make(map[string]any),
		store:  store,
	}
}

// Name returns the name of the session.
func (s *Session) Name() string {
	return s.name
}

// Set sets a key/value pair in the session state.
func (s *Session) Set(key string, value any) {
	s.values[key] = value
}

// Get returns the state value for the given key.
func (s *Session) Get(key string) any {
	return s.values[key]
}

// GetOk returns the state value for the given key and whether they key exists.
func (s *Session) GetOk(key string) (any, bool) {
	value, ok := s.values[key]
	return value, ok
}

// Save adds or updates the session. Identical to calling
// store.Save(w, session).
func (s *Session) Save(w http.ResponseWriter) error {
	return s.store.Save(w, s)
}

// Destroy destroys the session. Identical to calling
// store.Destroy(w, session.name).
func (s *Session) Destroy(w http.ResponseWriter) {
	s.store.Destroy(w, s.name)
}
