package main

import(
	"os"
	//"fmt"
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

	defs := []muxes.FuncDefinition{
		{"/", "^$", HelloWorld},
		{"/eo/", "^$", SalutonMondo},
		{"/test/", "", muxes.GetTest},
	}

	mux := muxes.DefineFuncs(nil, defs)
	srv := http.Server {
		Addr: *AddrStr,
		Handler: mux,
	}

	log.Printf("%s: running on '%s'\n",
		os.Args[0],
		*AddrStr)
	log.Fatal(srv.ListenAndServe())
}
