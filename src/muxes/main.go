package muxes

import(
	//"html/template"
	"net/http"
	"net/url"
	"regexp"
	//"strconv"
	//"io/ioutil"
	"fmt"
	"github.com/surdeus/ghost/src/urlpath"
)

type HndlArg struct {
	W http.ResponseWriter
	R *http.Request
	Q url.Values
	P string
}

type Handler func(a HndlArg)

type FuncDefinition struct {
	Pref, Re string
	Fn Handler
}

func MakeHttpHandleFunc(pref string, re *regexp.Regexp, fn Handler) http.HandlerFunc {
return func(w http.ResponseWriter, r *http.Request) {
	var(
		a HndlArg
		e error
	)

	a.P = r.URL.Path[len(pref):]
	if !urlpath.Validify(a.P, re) {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET" :
		a.Q, e = url.ParseQuery(r.URL.RawQuery)
	case "POST" :
		r.ParseForm()
	}

	if e != nil {
	}

	a.W = w
	a.R = r
	
	fn(a)
}}

func DefineFuncs(mux *http.ServeMux,defs []FuncDefinition) *http.ServeMux {
	if mux == nil {
		mux = http.NewServeMux()
	}

	for _, def := range defs {
		mux.HandleFunc(def.Pref,
			MakeHttpHandleFunc(def.Pref,
				regexp.MustCompile(def.Re),
				def.Fn))
	}

	return mux
}

func GetTest(a HndlArg){
	w := a.W
	r := a.R
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Path: '%s'\nRawQuery:'%s'\n", r.URL.Path, r.URL.RawQuery)
	fmt.Fprintf(w, "a.P: '%s'\n", a.P)
	fmt.Fprintf(w, "a.Q:\n")
	for k, v := range a.Q {
		fmt.Fprintf(w, "\t'%s':\n", k)
		for _, s := range v {
			fmt.Fprintf(w, "\t\t'%s'\n", s)
		}
	}

}


