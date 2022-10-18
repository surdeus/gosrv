package cookies

import (
	"net/http"
)

func ByName(cookies []*http.Cookie, name string) (string, bool) {
	for _, v := range cookies {
		if v.Name == name {
			return v.Value, true
		}
	}

	return "", false
}
