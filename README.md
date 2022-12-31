# sessions [![GoDoc](https://pkg.go.dev/badge/github.com/dghubble/sessions.svg)](https://pkg.go.dev/github.com/dghubble/sessions) [![Workflow](https://github.com/dghubble/sessions/actions/workflows/test.yaml/badge.svg)](https://github.com/dghubble/sessions/actions/workflows/test.yaml?query=branch%3Amain) [![Sponsors](https://img.shields.io/github/sponsors/dghubble?logo=github)](https://github.com/sponsors/dghubble) [![Mastodon](https://img.shields.io/badge/follow-news-6364ff?logo=mastodon)](https://fosstodon.org/@typhoon)

Package `sessions` provides minimalist Go sessions, backed by `securecookie` or database stores.

### Features

* `Store` provides an interface for managing sessions.
    * `New` returns a new named `Session`.
    * `Get` returns the named `Session` from the `http.Request` iff it was correctly verified and decoded. Otherwise the error is non-nil.
    * `Save` encodes and signs Session.Value data.
    * `Destroy` removes (expires) the session cookie of a given name.
* Each `Session` provides `Save` and `Destroy` convenience methods.
* Provides `CookieStore` for managing client-side secure cookies.
* Extensible for custom session database backends.

## Install

```
go get github.com/dghubble/sessions
```

## Documentation

Read [GoDoc](https://godoc.org/github.com/dghubble/sessions)

## Usage

Create a `Store` for managing `Session`'s. `NewCookieStore` returns a `Store` that signs and optionally encrypts cookies to support user sessions.

```go
import (
  "github.com/dghubble/sessions"
)

func NewServer() (http.Handler) {
  ...
  // client-side cookies
  sessionProvider := sessions.NewCookieStore(
    // use a 32 byte or 64 byte hash key
    []byte("signing-secret"),
    // use a 32 byte (AES-256) encryption key
    []byte("encryption-secret")
  )
  sessionProvider.Config.SameSite = http.SameSiteStrictMode
  ...
}
```

Issue a session cookie from a handler (e.g. login handler).

```go
func (s server) Login() http.Handler {
  fn := func(w http.ResponseWriter, req *http.Request) {
    // create a session
    session := s.sessions.New("my-app")
    // add user-id to session
    session.Set("user-id", 123)
    // save the session to the response
    if err := session.Save(w); err != nil {
      // handle error
    }
    ...
  }
  return http.HandlerFunc(fn)
}
```

Access the session and its values (e.g. require authentication).

```go
func (s server) RequireLogin() http.Handler {
  fn := func(w http.ResponseWriter, req *http.Request) {
    session, err := s.sessions.Get("my-app")
    if err != nil {
      http.Error(w, "missing session", http.StatusUnauthorized)
      return
    }

    userID := session.Get("user-id")
    fmt.Fprintf(w, `<p>Welcome %d!</p>
    <form action="/logout" method="post">
    <input type="submit" value="Logout">
    </form>`, userID)
  }
  return http.HandlerFunc(fn)
}
```

Delete a session when a user logs out.

```go
func (s server) Logout() http.Handler {
  fn := func(w http.ResponseWriter, req *http.Request) {
    if req.Method == "POST" {
      s.sessions.Destroy(w, "my-app")
    }
    http.Redirect(w, req, "/", http.StatusFound)
  }
  return http.HandlerFunc(fn)
}
```

### Differences from gorilla/sessions

* Gorilla stores a context map of Requests to Sessions to abstract multiple sessions. `dghubble/sessions` provides individual sessions, leaving multiple sessions to a `multisessions` package. No Registry is needed.
* Gorilla has a depedency on `gorilla/context`, a non-standard context.
* Gorilla requires all handlers be wrapped in `context.ClearHandler` to avoid memory leaks.
* Gorilla's `Store` interface is surprising. `New` and `Get` can both possibly return a new session, a field check is needed. Some use cases expect developers to [ignore an error](https://github.com/gorilla/sessions/blob/master/doc.go#L32). `Destroy` isn't provided.

## License

[MIT License](LICENSE)
