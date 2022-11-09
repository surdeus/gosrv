package cookiex

import (
	"net/http"
	"time"
)

func ByName(cookies []*http.Cookie, name string) (string, bool) {
	for _, v := range cookies {
		if v.Name == name {
			return v.Value, true
		}
	}

	return "", false
}

func Delete(w http.ResponseWriter, k string) {
	c := &http.Cookie {
		Name: k,
		Value: "",
		Path: "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, c)
}

