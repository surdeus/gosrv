package rex

import(
	"regexp"
)

func
Validify(u string, re *regexp.Regexp) bool {
	if re == nil {
		return true
	}
	
	ret := re.Find([]byte(u))
	if ret == nil {
		return false
	}

	return true
}

