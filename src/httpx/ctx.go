package httpx

import (
	"net/http"
	"net/url"
	"time"
)

// The type describes structure
// to save current state of the handler,
// make handlers take less arguments
// and is supposed to be used in chaining.
// The fields are public for more flexibility.
type Context struct {
	// The standard http writer.	
	W http.ResponseWriter
	// Request.
	R *http.Request
	// Query values.
	Q url.Values
	// Part of path without prefix.
	P string
	// To save values between chained hanlders.
	V map[string] any
}

// Makes cookie to expire on the browser side.
func (a *Context) DeleteCookie(k string) {
	c := &http.Cookie {
		Name: k,
		Value: "",
		Path: "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(a.W, c)
}

func (a *Context) Cookies() Cookies {
	return a.R.Cookies()
}

// Sends to the writer default Golang "Not Found".
func (a *Context) NotFound() {
	http.NotFound(a.W, a.R)
}

// Return values of parsed form as a map.
func (a *Context) Values() url.Values {
	return a.R.URL.Query()
}

// Sends the internal server error message with
// custom error.
func (a *Context) ServerError(err error) {
	http.Error(a.W, err.Error(), http.StatusInternalServerError)
}