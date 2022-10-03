package main

import(
	"os"
	"fmt"
	"flag"
	"log"
	"net/http"
	"github.com/surdeus/ghost/src/router"
	//"regexp"
)

func HelloWorld(w http.ResponseWriter, r *http.Request,
		a router.HndlArg) {
	fmt.Fprintf(w, "Hello, World!")
}

func SalutonMondo(w http.ResponseWriter, r *http.Request,
		a router.HndlArg) {
	name := "Mondo"
	_, ok := a.Q["name"]
	if ok {
		name = a.Q["name"][0]
	}
	fmt.Fprintf(w, "Saluton, %s!", name)
}

func
main(){
	AddrStr := flag.String("a", ":8080", "Adress string")
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		os.Exit(1)
	}

	defs := []router.Definition{
		{"/", "^$", HelloWorld},
		{"/eo/", "^$", SalutonMondo},
		{"/test/", "", router.GetTest},
	}

	mux := router.Mux(nil, defs)
	srv := http.Server {
		Addr: *AddrStr,
		Handler: mux,
	}

	log.Printf("%s: running on '%s'\n",
		os.Args[0],
		*AddrStr)
	log.Fatal(srv.ListenAndServe())
}
