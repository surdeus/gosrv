package main

import(
	"os"
	"fmt"
	"flag"
	"log"
	"net/http"
	"github.com/surdeus/ghost/src/muxes"
	"github.com/surdeus/ghost/src/templates"
	//"regexp"
)

var (
	tmpls *templates.Templates
)

func HelloWorld(a muxes.HndlArg) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpls.Exec(a.W, "hellos/en", struct{}{})
}

func SalutonMondo(a muxes.HndlArg) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	name := "Mondo"
	_, ok := a.Q["name"]
	if ok {
		name = a.Q["name"][0]
	}
	tmpls.Exec(a.W, "hellos/eo", struct{Name string}{Name: name})
}

func GeneralChainFunc(hndl muxes.Handler) muxes.Handler {
return func(a muxes.HndlArg) {
	fmt.Println("general function got called")
	hndl(a)
}}

func OtherFunc(hndl muxes.Handler) muxes.Handler {
return func(a muxes.HndlArg) {
	fmt.Println("some other func called")
	hndl(a)
}}

func main(){
	AddrStr := flag.String("a", ":8080", "Adress string")
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		os.Exit(1)
	}

	cfg := templates.ParseConfig{
		Root: "tmpl",
		FuncMap: templates.FuncMap{
			"SomeFunc": func() string {
				return "<div>This is some string</div>"
		}},
	}

	var err error
	tmpls, err = templates.Parse(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", tmpls)

	defs := []muxes.HndlDef {
		{"/", "^$", muxes.Handlers{
			"GET": muxes.Chained(muxes.Chain{GeneralChainFunc}, HelloWorld)} },
		{"/eo/", "^$", muxes.Handlers{
			"GET":muxes.Chained(muxes.Chain{GeneralChainFunc, OtherFunc}, SalutonMondo)} },
		{"/test/", "", muxes.Handlers{"GET": muxes.GetTest} },
	}

	mux := muxes.Define(nil, defs)
	srv := http.Server {
		Addr: *AddrStr,
		Handler: mux,
	}

	log.Printf("%s: running on '%s'\n",
		os.Args[0],
		*AddrStr)
	log.Fatal(srv.ListenAndServe())
}
