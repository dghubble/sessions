# sessions Changelog

Notable changes between releases.

## Latest

* Allow `Session` to store values with specified type (`V`) (i.e. generics) ([#22](https://github.com/dghubble/sessions/pull/21))
  * `Session` state is now a `map[string]V` instead of a `map[string]any`
  * Update `Set`, `Get`, and `GetOk` methods to use generic type `V`
  * Change `Session` to `Session[V any]` to specify the type of value stored in the Session
  * See updated usage docs for examples
* Change `Store` to `Store[V any]` to specify the type of value stored in sessions
* Change `NewCookieStore` to `NewCookieStore[V any]` to specify the type of value stored in sessions

## v0.3.0

* Change `CookieStore` and its fields to be non-exported ([#19](https://github.com/dghubble/sessions/pull/19))
  * Change `NewCookieStore` to require a `*CookieConfig` and return a `Store`
  * Rename `Config` struct to `CookieConfig`
  * Add `DefaultCookieConfig` and `DebugCookieConfig` convenience variables
* Change the `Session` field `Values` to be non-exported ([#18](https://github.com/dghubble/sessions/pull/18))
  * Add `Session` `Set` method to set a key/value pair
  * Add `Session` `Get` method to get a value for a given key
  * Add `Session` `GetOk` to get a value for a given key and whether the key exists in the map
* Remove cookie `Config` field from `Session` ([#17](https://github.com/dghubble/sessions/pull/17))

## v0.2.1

* Update minimum Go version from v1.17 to v1.18 ([#15](https://github.com/dghubble/sessions/pull/15))

## v0.2.0

* Fix `go.mod` to include `gorilla/securecookie` ([#7](https://github.com/dghubble/sessions/pull/7))

## v0.1.0

* Initial release
* Require Go v1.11+
