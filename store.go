package sessions

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

// A Store manages creating, accessing, writing, and expiring Sessions.
type Store[V any] interface {
	// New returns a new named Session
	New(name string) *Session[V]
	// Get a named Session from the request
	Get(req *http.Request, name string) (*Session[V], error)
	// Save writes a Session to the ResponseWriter
	Save(w http.ResponseWriter, session *Session[V]) error
	// Destroy removes (expires) a named Session
	Destroy(w http.ResponseWriter, name string)
}

var _ Store[any] = &cookieStore[any]{}

// CookieStore stores Sessions in secure cookies (i.e. client-side)
type cookieStore[V any] struct {
	config *CookieConfig
	// encodes and decodes signed and optionally encrypted cookie values
	codecs []securecookie.Codec
}

// NewCookieStore returns a new Store that signs and optionally encrypts
// session state in http cookies.
func NewCookieStore[V any](config *CookieConfig, keyPairs ...[]byte) Store[V] {
	if config == nil {
		config = DefaultCookieConfig
	}

	return &cookieStore[V]{
		config: config,
		codecs: securecookie.CodecsFromPairs(keyPairs...),
	}
}

// New returns a new named Session.
func (s *cookieStore[V]) New(name string) *Session[V] {
	return NewSession[V](s, name)
}

// Get returns the named Session from the Request. Returns an error if the
// session cookie cannot be found, the cookie verification fails, or an error
// occurs decoding the cookie value.
func (s *cookieStore[V]) Get(req *http.Request, name string) (session *Session[V], err error) {
	cookie, err := req.Cookie(name)
	if err == nil {
		session = s.New(name)
		err = securecookie.DecodeMulti(name, cookie.Value, &session.values, s.codecs...)
	}
	return session, err
}

// Save adds or updates the Session on the response via a signed and optionally
// encrypted session cookie. Session Values are encoded into the cookie value
// and the session Config sets cookie properties.
func (s *cookieStore[V]) Save(w http.ResponseWriter, session *Session[V]) error {
	cookieValue, err := securecookie.EncodeMulti(session.Name(), &session.values, s.codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, newCookie(session.Name(), cookieValue, s.config))
	return nil
}

// Destroy deletes the Session with the given name by issuing an expired
// session cookie with the same name.
func (s *cookieStore[V]) Destroy(w http.ResponseWriter, name string) {
	http.SetCookie(w, newCookie(name, "", &CookieConfig{MaxAge: -1, Path: s.config.Path}))
}
