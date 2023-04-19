package httpx

type ChainHandler func(h HandlerFunc) HandlerFunc
type Chain []ChainHandler

// Chain functions into final form.
func Chained(c Chain, h HandlerFunc) HandlerFunc {
	if len(c) > 1 {
		return c[0]( Chained(c[1:], h) )
	}
	
	if len(c) == 0 {
		return h
	}

	return c[0](h)
}
