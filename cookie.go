package sessions

import (
	"net/http"
	"time"
)

// Config is the set of cookie properties.
type Config struct {
	// cookie domain/path scope (leave zeroed for requested resource scope)
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge int
	// cookie may only be transferred over HTTPS
	Secure bool
	// browser should prohibit non-HTTP (i.e. javascript) cookie access
	HTTPOnly bool
	// prohibit sending in cross-site requests with SameSiteLaxMode or SameSiteLaxMode
	SameSite http.SameSite
}

// newCookie returns a new http.Cookie with the given name, value, and
// properties from config.
func newCookie(name, value string, config *Config) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   config.MaxAge,
		HttpOnly: config.HTTPOnly,
		Secure:   config.Secure,
		SameSite: config.SameSite,
	}
	// IE <9 does not understand MaxAge, set Expires based on MaxAge
	if expires, present := cookieExpires(config.MaxAge); present {
		cookie.Expires = expires
	}
	return cookie
}

// cookieExpires takes the MaxAge number of seconds a Cookie should be valid
// and returns the Expires time.Time and whether the attribtue should be set.
// http://golang.org/src/net/http/cookie.go?s=618:801#L23
func cookieExpires(maxAge int) (time.Time, bool) {
	if maxAge > 0 {
		d := time.Duration(maxAge) * time.Second
		return time.Now().Add(d), true
	} else if maxAge < 0 {
		return time.Unix(1, 0), true // first second of the epoch
	}
	return time.Time{}, false
}
