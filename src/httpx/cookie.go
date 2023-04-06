package httpx

import (
	"net/http"
	//"time"
)

type Cookie = http.Cookie
type Cookies []*Cookie

func (cookies Cookies) Get(name string) (string, bool) {
	for _, v := range cookies {
		if v.Name == name {
			return v.Value, true
		}
	}

	return "", false
}
