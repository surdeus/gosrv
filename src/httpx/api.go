package httpx

import (
	"encoding/gob"
)

// Implementation of API interfaces for
// simpler client-server interaction
// with default HTTP requests.

// Context for API functions to be able
// to read since it requires concrete types
// to be know so it is done in the handling
// functions.
type ApiContext struct {
	*Context
	dec *gob.Decoder
}

// Function that implements API.
type ApiHandlerFunc func(*ApiContext) (chan any, error)

// Wraps new handler from API handler.
func makeApiHandler(h ApiHandlerFunc) HandlerFunc {
return func(c *Context) {
	apiContext := &ApiContext{
		Context: c,
	}
	
	apiContext.dec = gob.NewDecoder(c.R.Body)
	h(apiContext)
}}

func (c *ApiContext)ReadValues(ptrs ...any) (int, error) {
	for i, ptr := range ptrs {
		err := c.dec.Decode(ptr)
		if err != nil {
			return i, err
		}
	}
	
	return -1, nil
}

