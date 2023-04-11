package httpx

import(
	"net/http"
	"net/url"
	"log"
	"fmt"
	"strings"
	"github.com/surdeus/gosrv/src/rex"
)

type Server = http.Server

type ResponseWriter = http.ResponseWriter
type Request = http.Request

type Mux struct {
	*http.ServeMux
}

type HandlerFunc func(a *Context)
type Handler interface {
	ServeHttp(*Context)
}

type HandlerFuncMap map[string] HandlerFunc

// Create final function handler.
func makeHttpHandleFunc(
	def HandlerDef,
) http.HandlerFunc {
	if def.simpleHndl != nil {
		return def.simpleHndl
	} else if def.simpleHndler != nil {
		return def.simpleHndler.ServeHTTP
	}
	
	pref := def.pref
	re := def.re
	handlers := def.handlers
	
return func(w http.ResponseWriter, r *http.Request) {
	
	var(
		a Context
		e error
	)

	p := r.URL.Path
	if p == ""  || p[len(p)-1] != '/' {
		p += "/"
	}
	a.P = p[len(pref):]
	
	if re != nil && !rex.Validify(a.P, re) {
		http.NotFound(w, r)
		return
	}

	// Parsing of arguments and shit.
	method := r.Method
	
	method = strings.ToUpper(method)	
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
	handler(&a)
}}

// Returns new empty mux.
func NewMux() *Mux {
	return &Mux{
		http.NewServeMux(),
	}
}

// Define new handlers by HandlerDef structure. 
func (mux *Mux)Define(defs ...HandlerDef) {
	for _, def := range defs {
		hndl := makeHttpHandleFunc(def)
		mux.HandleFunc(def.pref, hndl)
	}
}

// Simple check function for debug.
func GetTest(a *Context){
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


