package httpx

import (
	"encoding/gob"
	"net/http"
	"bytes"
	"errors"
	"io"
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
	gobDec *gob.Decoder
	err error
	i int
}

// Function that implements API.
type ApiHandlerFunc func(*ApiContext) (chan any, error)

type ApiResponseType int

type ApiResponse struct {
	dec *gob.Decoder
	resp *http.Response
	err error
	i int
}

const (
	ApiResponseTypeNone ApiResponseType = iota
	ApiResponseTypeError
	ApiResponseTypeSuccess
	ApiResponseTypeLast
)

var (
	WrongResponseTypeErr = errors.New("wrong response type")
)

// Client function to make a GOB query.
// Cannot do JSON queries, since this is a Go client.
func ApiQuery(u string, vals ...any) (*ApiResponse, error) {
	/* Note: should implement channel handling so
		so we can send many values in parallel. */
		
	bts := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(bts)
	
	for _, val := range vals {
		err := enc.Encode(val)
		if err != nil {
			return nil, err
		}
	}
	
	resp, err := http.Post(
		u, "application/gob", bts,
	)
	
	if err != nil {
		return nil, err
	}
	
	dec := gob.NewDecoder(resp.Body)
	
	// Checking if API returns error.
	var respType ApiResponseType
	dec.Decode(&respType)
	switch respType {
	case ApiResponseTypeError :
		var err error
		dec.Decode(&err)
		return nil, err
	case ApiResponseTypeSuccess :
		return &ApiResponse{
			dec: dec,
			resp: resp,
		}, nil
	default :
		return nil, WrongResponseTypeErr
	}
	
}

func (resp *ApiResponse) Done() bool {
	return resp.err != nil
}

func (resp *ApiResponse) Scan(ptrs ...any) bool {
	for i, v := range ptrs {
		err := resp.dec.Decode(v)
		if err != nil {
			if err != io.EOF {
				resp.err = err
			}
			resp.i = i
			return false
		}
	}
	
	return true
}

func (resp *ApiResponse) Err() error {
	return resp.err
}

func (resp *ApiResponse) ErrI() int {
	return resp.i
}

// Wraps new handler from API handler.
func makeApiHandler(h ApiHandlerFunc) HandlerFunc {
return func(c *Context) {
	apiContext := &ApiContext{
		Context: c,
	}
	
	apiContext.gobDec = gob.NewDecoder(c.R.Body)
	chn, err := h(apiContext)
	enc := gob.NewEncoder(c.W)
	
	c.W.Header().Set("Content-type", "application/gob")
	if err != nil {
		enc.Encode(ApiResponseTypeError)
		enc.Encode(err)
		return
	}
	
	enc.Encode(ApiResponseTypeSuccess)
	for v := range chn {
		err := enc.Encode(v)
		// If we get an error here
		// a programmer implementing API
		// is doing something wrong so we panic,
		// since this is writing from the server it
		// should be implemented correctly.
		if err != nil {
			panic(err)
		}
	}
}}

// Reads transported values into buffers that pointers point to.
// Returns index of argument where it gets error.
func (c *ApiContext) Scan(ptrs ...any) (bool) {
	for i, ptr := range ptrs {
		err := c.gobDec.Decode(ptr)
		if err != nil {
			if err != io.EOF {
				c.err = err
			}
			c.i = i
			return false
		}
	}
	
	return true
}

func (c *ApiContext) Err() error {
	return c.err
}

func (c *ApiContext) ErrI() int {
	return c.i
}

