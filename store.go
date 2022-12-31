package sessions

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

// Store is the interface for creating, reading, updating and destroying
// named Sessions.
type Store interface {
	New(name string) *Session
	Get(req *http.Request, name string) (*Session, error)
	Save(w http.ResponseWriter, session *Session) error
	Destroy(w http.ResponseWriter, name string)
}

// CookieStore stores Sessions in secure cookies (i.e. client-side)
type CookieStore struct {
	// encodes and decodes signed and optionally encrypted cookie values
	Codecs []securecookie.Codec
	// configures session cookie properties of new Sessions
	Config *Config
}

// NewCookieStore returns a new CookieStore which signs and optionally encrypts
// session cookies.
func NewCookieStore(keyPairs ...[]byte) *CookieStore {
	return &CookieStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Config: &Config{
			Path:     "/",
			MaxAge:   defaultMaxAge,
			HTTPOnly: true,
			SameSite: http.SameSiteDefaultMode,
		},
	}
}

// New returns a new Session with the requested name and the store's config
// value.
func (s *CookieStore) New(name string) *Session {
	return NewSession(s, name)
}

// Get returns the named Session from the Request. Returns an error if the
// session cookie cannot be found, the cookie verification fails, or an error
// occurs decoding the cookie value.
func (s *CookieStore) Get(req *http.Request, name string) (session *Session, err error) {
	cookie, err := req.Cookie(name)
	if err == nil {
		session = s.New(name)
		err = securecookie.DecodeMulti(name, cookie.Value, &session.Values, s.Codecs...)
	}
	return session, err
}

// Save adds or updates the Session on the response via a signed and optionally
// encrypted session cookie. Session Values are encoded into the cookie value
// and the session Config sets cookie properties.
func (s *CookieStore) Save(w http.ResponseWriter, session *Session) error {
	cookieValue, err := securecookie.EncodeMulti(session.Name(), &session.Values, s.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, newCookie(session.Name(), cookieValue, s.Config))
	return nil
}

// Destroy deletes the Session with the given name by issuing an expired
// session cookie with the same name.
func (s *CookieStore) Destroy(w http.ResponseWriter, name string) {
	http.SetCookie(w, newCookie(name, "", &Config{MaxAge: -1, Path: s.Config.Path}))
}
