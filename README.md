# sessions
[![GoDoc](https://pkg.go.dev/badge/github.com/dghubble/sessions.svg)](https://pkg.go.dev/github.com/dghubble/sessions)
[![Workflow](https://github.com/dghubble/sessions/actions/workflows/test.yaml/badge.svg)](https://github.com/dghubble/sessions/actions/workflows/test.yaml?query=branch%3Amain)
[![Sponsors](https://img.shields.io/github/sponsors/dghubble?logo=github)](https://github.com/sponsors/dghubble)
[![Mastodon](https://img.shields.io/badge/follow-news-6364ff?logo=mastodon)](https://fosstodon.org/@dghubble)

<img align="right" src="https://storage.googleapis.com/dghubble/small-gopher-with-cookie.png">

Package `sessions` provides minimalist Go sessions, backed by `securecookie` or database stores.

### Features

* `Store` provides an interface for managing a user `Session`
    * May be implemented by custom session database backends
* `Session` stores a typed value (via generics)
* `Session` provides convenient key/value `Set`, `Get`, and `GetOk` methods
* `NewCookieStore` implements a `Store` backed by client-side cookies (signed and optionally encrypted)

## Docs

Read [GoDoc](https://godoc.org/github.com/dghubble/sessions)

## Usage

Create a `Store` for managing `Session`'s. `NewCookieStore` returns a `Store` that signs and optionally encrypts cookies to support user sessions.

A `Session` stores a map of key/value pairs (e.g. "userID": "a1b2c3"). Starting with v0.4.0, `sessions` uses Go generics to allow specifying a type for stored values. Previously, values were type `interface{}` or `any`, which required type assertions.

```go
import (
  "github.com/dghubble/sessions"
)

func NewServer() (http.Handler) {
  ...
  // client-side cookies
  store := sessions.NewCookieStore[string](
    sessions.DefaultCookieConfig,
    // use a 32 byte or 64 byte hash key
    []byte("signing-secret"),
    // use a 32 byte (AES-256) encryption key
    []byte("encryption-secret")
  )
  ...
  server.sessions = store
}
```

Issue a session cookie from a handler (e.g. login handler).

```go
func (s server) Login() http.Handler {
  fn := func(w http.ResponseWriter, req *http.Request) {
    // create a session
    session := s.sessions.New("my-app")
    // add user-id to session
    session.Set("user-id", "a1b2c3")
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

    userID, present := session.GetOk("user-id")
    if !present {
      http.Error(w, "missing user-id", http.StatusUnauthorized)
      return
    }

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

## License

[MIT License](LICENSE)
