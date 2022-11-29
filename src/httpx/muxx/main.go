package muxx

import(
	//"html/template"
	"net/http"
	"net/url"
	"regexp"
	//"strconv"
	//"io/ioutil"
	"log"
	"fmt"
	"github.com/surdeus/gosrv/src/rex"
)

type Handler func(a HndlArg)
type ChainHandler func(h Handler) Handler
type Chain []ChainHandler
type Handlers map[string] Handler

type HndlArg struct {
	W http.ResponseWriter
	R *http.Request
	Q url.Values
	P string
	// Simpler than context.
	// Use contexts for more complex stuff.
	V map[string] any
}

type HndlDef struct {
	Pref, Re string
	Handlers Handlers
}

// Chain functions into final form.
func Chained(c Chain, h Handler) Handler {
	if len(c) > 1 {
		return c[0]( Chained(c[1:], h) )
	}

	return c[0](h)
}

func (a HndlArg)NotFound() {
	http.NotFound(a.W, a.R)
}

func (a HndlArg)Values() url.Values {
	return a.R.URL.Query()
}

func (a HndlArg)ServerError(err error) {
	http.Error(a.W, err.Error(), http.StatusInternalServerError)
}

// Create final function handler.
func MakeHttpHandleFunc(
	pref string,
	re *regexp.Regexp,
	handlers Handlers,
) http.HandlerFunc {
return func(w http.ResponseWriter, r *http.Request) {
	var(
		a HndlArg
		e error
	)

	p := r.URL.Path
	if p == ""  || p[len(p)-1] != '/' {
		p += "/"
	}
	a.P = p[len(pref):]
	if !rex.Validify(a.P, re) {
		http.NotFound(w, r)
		return
	}

	// Parsing of arguments and shit.
	method := r.Method
	switch method {
	case "GET" :
		a.Q, e = url.ParseQuery(r.URL.RawQuery)
	case "POST" :
		fallthrough
	case "PUT" :
		fallthrough
	case "PATCH" :
		r.ParseForm()
	}

	if e != nil {
		log.Println(e)
	}

	a.W = w
	a.R = r
	
	handler, ok := handlers[method]
	if !ok {
		handler, ok = handlers[""]
		if !ok {
			http.NotFound(w, r)
			return
		}
	}

	a.V = make(map[string] any, 5)
	handler(a)
}}

func Define(mux *http.ServeMux, defs []HndlDef) *http.ServeMux {
	if mux == nil {
		mux = http.NewServeMux()
	}

	for _, def := range defs {
		mux.HandleFunc(def.Pref,
			MakeHttpHandleFunc(def.Pref,
				regexp.MustCompile(def.Re),
				def.Handlers))
	}

	return mux
}

func DefineSimple(mux *http.ServeMux, pref string,
	hndl http.HandlerFunc) *http.ServeMux {
		if mux == nil {
			mux = http.NewServeMux()
		}
		mux.HandleFunc(pref, hndl)

		return mux
}

func DefineStatic(mux *http.ServeMux, path, pref string) *http.ServeMux {
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
}

func GetTest(a HndlArg){
	w := a.W
	r := a.R
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Path: '%s'\nRawQuery:'%s'\n", r.URL.Path, r.URL.RawQuery)
	fmt.Fprintf(w, "a.P: '%s'\n", a.P)
	fmt.Fprintln(w, "a.Q:\n")
	for k, v := range a.R.URL.Query() {
		fmt.Fprintf(w, "\t'%s':\n", k)
		for _, s := range v {
			fmt.Fprintf(w, "\t\t'%s'\n", s)
		}
	}

	fmt.Fprintf(w, "a.R.Cookies():\n")
	for _, c := range a.R.Cookies() {
		fmt.Fprintf(w, "\t'%v'\n", c)
	}
}


