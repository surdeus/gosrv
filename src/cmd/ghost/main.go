package main

import(
	"os"
	"flag"
	"log"
	"net/http"
	"github.com/mojosa-software/gosrv/src/httpx"
)


func main(){
	AddrStr := flag.String("a", ":8080", "Adress string")
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		os.Exit(1)
	}

	var pth string
	if len(args) == 1 {
		pth = args[0]
	} else {
		pth = "."
	}

	mux := httpx.NewMux()
	mux.Define(
		httpx.Def("/").StaticFiles(pth),
	)
	srv := http.Server {
		Addr: *AddrStr,
		Handler: mux,
	}
	log.Printf("%s: Trying to run on '%s'...\n",
		os.Args[0],
		*AddrStr)
	log.Fatal(srv.ListenAndServe())
}

