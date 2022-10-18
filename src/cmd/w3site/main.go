package main

import(
	"os"
	"fmt"
	"flag"
	"log"
	"net/http"
	"github.com/surdeus/ghost/src/muxes"
	"github.com/surdeus/ghost/src/templates"
	//"html/template"
	//"regexp"
)

var (
	tmpls templates.Templates
)

func HelloWorld(a muxes.HndlArg) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := "default"
	_, ok := a.Q["tmpl"]
	if ok {
		tmpl = a.Q["tmpl"][0]
	}
	tmpls.Exec(a.W, tmpl, "hellos/en", struct{}{})
}

func SalutonMondo(a muxes.HndlArg) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	name := "Mondo"
	_, ok := a.Q["name"]
	if ok {
		name = a.Q["name"][0]
	}
	tmpls.Exec(a.W, "default", "hellos/eo", struct{Name string}{Name: name})
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
		Component: "tmpl/c/",
		View: "tmpl/v/",
		Template: "tmpl/t/",
		FuncMap: templates.FuncMap{
			"SomeFunc": func() string {
				return "<div>This is some string</div>"
			},
			/*"TmplFunc" : func(template *Template) {
				fmt.Printf("got '%s template\n',")
			},*/
		},
	}

	var err error
	tmpls, err = templates.Parse(cfg)
	if err != nil {
		panic(err)
	}


	/*fmt.Println("Parsed templates:")
	for _, v := range tmpls.Templates() {
		fmt.Printf("'%s'\n", v.Name())
	}*/

	fmt.Printf("%v\n", tmpls)

	defs := []muxes.HndlDef {
		{"/", "^$", muxes.Handlers{
			"GET": muxes.Chained(muxes.Chain{GeneralChainFunc}, HelloWorld)} },
		{"/eo/", "^$", muxes.Handlers{
			"GET":muxes.Chained(muxes.Chain{GeneralChainFunc, OtherFunc}, SalutonMondo)} },
		{"/test/", "", muxes.Handlers{"GET": muxes.GetTest} },
	}

	mux := muxes.Define(nil, defs)
	muxes.DefineStatic(mux, "s/", "/s/")
	srv := http.Server {
		Addr: *AddrStr,
		Handler: mux,
	}

	log.Printf("%s: running on '%s'\n",
		os.Args[0],
		*AddrStr)
	log.Fatal(srv.ListenAndServe())
}
