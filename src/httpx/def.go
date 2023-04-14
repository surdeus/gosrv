package httpx

import (
	"net/http"
	"strings"
	"regexp"
)

type HandlerType int

type HandlerDef struct {
	typ HandlerType
	pref string
	re *regexp.Regexp
	handlers HandlerFuncMap
	simpleHandlerFunc http.HandlerFunc
	simpleHandler http.Handler
}

const (
	HandlerTypeNone = iota
	HandlerTypeDefault 
	HandlerTypeSimpleFunc
	HandlerTypeSimple
)

func Def(pref string) HandlerDef {
	return HandlerDef{
		pref: pref,
		re: regexp.MustCompile(""),
		handlers: make(HandlerFuncMap),
	}
}

func (d HandlerDef)Re(re string) HandlerDef {
	d.re = regexp.MustCompile(re)
	return d
}

func (d HandlerDef) Method(m string, h HandlerFunc) HandlerDef {
	d.typ = HandlerTypeDefault
	d.handlers[strings.ToUpper(m)] = h
	return d
}

func (d HandlerDef) Api(h ApiHandlerFunc) HandlerDef {
	d.typ = HandlerTypeDefault
	
	// Clearing map so we have no collisions.
	d.handlers = make(HandlerFuncMap)
	d.handlers[MethodEmpty] = makeApiHandler(h)
	
	return d
}

func (d HandlerDef) SimpleHandlerFunc(s http.HandlerFunc) HandlerDef {
	d.typ = HandlerTypeSimpleFunc
	d.simpleHandlerFunc = s
	
	return d
}

func (d HandlerDef) SimpleHandler(s http.Handler) HandlerDef {
	d.typ = HandlerTypeSimple
	d.simpleHandler = s
	
	return d
}

func (d HandlerDef) StaticFiles(path string) HandlerDef {
	d.typ = HandlerTypeSimple
	fs := http.FileServer(http.Dir(path))
	d.simpleHandler = http.StripPrefix(d.pref, fs)
	
	return d
}

