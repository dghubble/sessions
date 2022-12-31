package sessions

import (
	"net/http"
)

const (
	defaultMaxAge = 3600 * 24 * 7 // 1 week
)

// Session represents Values state which  a named bundle of maintained web state
// stores web session state
type Session struct {
	name   string // session cookie name
	Values map[string]interface{}
	// convenience methods Save and Destroy use store
	store Store
}

// NewSession returns a new Session.
func NewSession(store Store, name string) *Session {
	return &Session{
		store:  store,
		name:   name,
		Values: make(map[string]interface{}),
	}
}

// Name returns the name of the session.
func (s *Session) Name() string {
	return s.name
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
