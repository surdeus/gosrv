package httpx

import (
	"net/http"
	"strings"
	"regexp"
)

type HandlerDef struct {
	pref string
	re *regexp.Regexp
	handlers HandlerFuncMap
	simpleHndl http.HandlerFunc
	simpleHndler http.Handler
}

func Def(pref string) HandlerDef {
	return HandlerDef{
		pref: pref,
		re: nil,
		handlers: make(HandlerFuncMap),
	}
}

func (d HandlerDef)Re(re string) HandlerDef {
	d.re = regexp.MustCompile(re)
	return d
}

func (d HandlerDef) Method(m string, h HandlerFunc) HandlerDef {
	d.handlers[strings.ToUpper(m)] = h
	return d
}

func (d HandlerDef) SimpleHandlerFunc(s http.HandlerFunc) HandlerDef {
	d.simpleHndler = nil
	d.simpleHndl = s
	return d
}

func (d HandlerDef) SimpleHandler(s http.Handler) HandlerDef {
	d.simpleHndl = nil
	d.simpleHndler = s
	return d
}

func (d HandlerDef) StaticFiles(path string) HandlerDef {
	fs := http.FileServer(http.Dir(path))
	d.simpleHndler = http.StripPrefix(d.pref, fs)
	
	return d
}

/*func DefineStatic(mux *http.ServeMux, path, pref string) *http.ServeMux {
	if mux == nil {
		mux = http.NewServeMux()
	}

	fs := http.FileServer(http.Dir(path))
	mux.Handle(
		pref,
		http.StripPrefix(
				pref,
				fs,	
		),
	)

	return mux
}*/
