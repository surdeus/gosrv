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
	tmpls templates.Templates
)

func HelloWorld(w http.ResponseWriter, r *http.Request,
		a muxes.HndlArg) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpls.Execute(w, "hello", struct{}{})
}

func SalutonMondo(w http.ResponseWriter, r *http.Request,
		a muxes.HndlArg) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	name := "Mondo"
	_, ok := a.Q["name"]
	if ok {
		name = a.Q["name"][0]
	}
	tmpls.Execute(w, "saluton", struct{Name string}{Name: name})
}

func main(){
	AddrStr := flag.String("a", ":8080", "Adress string")
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		os.Exit(1)
	}

	cfg := templates.ParseConfig{
		Gen: "tmpl/gen",
		Sep: "tmpl/sep",
		FuncMap: templates.FuncMap{
			"SomeFunc": func() string {
				return "<div>This is some string</div>"
		}},
	}

	tmpls = templates.MustParseTemplates(cfg)

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
