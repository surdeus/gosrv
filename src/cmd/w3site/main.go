package main

import(
	"os"
	"fmt"
	"flag"
	"log"
	"net/http"
	"encoding/json"
	"math/rand"
	"encoding/base64"
	"github.com/surdeus/ghost/src/muxes"
	"github.com/surdeus/ghost/src/templates"
	"github.com/surdeus/ghost/src/cookies"
	//"html/template"
	//"regexp"
)

type Users map[string] string

var (
	tokens = make(map[string] string)
	tmpls templates.Templates
	users Users
	datPath = "dat/"
	dbPath = datPath+"db/"
	staticPath = datPath+"s/"
	usersDbPath = dbPath + "users"
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

func GetCookies(a muxes.HndlArg) {
	_, ok1 := a.Q["name"]
	_, ok2 := a.Q["value"]
	if !ok1 || !ok2 {
		http.Error(a.W, "Wrong args", http.StatusInternalServerError)
		return
	}


	cookie := &http.Cookie{
		Name: a.Q["name"][0],
		Value: a.Q["value"][0],
		Path: "/",
	}
	http.SetCookie(a.W, cookie)
	a.W.WriteHeader(200)
	a.W.Write([]byte("success"))
}

func LoginGet(a muxes.HndlArg) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpls.Exec(a.W, "unauth", "login", nil)
}

func LoginPost(a muxes.HndlArg) {
	var (
		token string
	)
	formEmail := a.R.Form.Get("email")
	formPassword := a.R.Form.Get("password")
	fmt.Printf(
		"Email: '%s', password: '%s'\n",
		formEmail,
		formPassword,
	)

	password, ok := users[formEmail]
	if !ok {
		http.NotFound(a.W, a.R)
		return
	} else if password != formPassword  {
		a.W.WriteHeader(http.StatusUnauthorized)
		a.W.Write([]byte("Unauth"))
		return
	}

	token, _ = GenerateToken(formEmail, formPassword)
	tokens[token] = formEmail

	cookie := &http.Cookie{
		Name: "auth-token",
		Value: token,
		Path: "/",
	}

	http.SetCookie(a.W, cookie)
	http.Redirect(a.W, a.R, "/", http.StatusFound)
}

func GenerateToken(email, password string) (string, error) {
	token := make([]byte, 256)
	rand.Read(token)
	return base64.StdEncoding.EncodeToString(token), nil
}

func Authorize(hndl muxes.Handler) muxes.Handler {
return func(a muxes.HndlArg) {
	cookie := a.R.Cookies()

	authToken, ok := cookies.ByName(cookie, "auth-token")
	if !ok {
		http.Redirect(a.W, a.R,
			"/login/",
			http.StatusFound,
		)
		return
	} 

	_, ok = tokens[authToken]
	if !ok {
		http.Redirect(a.W, a.R,
			"/login/",
			http.StatusUnauthorized,
		)
	}

	hndl(a)
}}

func Unauthorize(hndl muxes.Handler) muxes.Handler {
return func(a muxes.HndlArg) {
	cookie := a.R.Cookies()

	_, ok := cookies.ByName(cookie, "auth-token")
	if ok {
		http.Redirect(a.W, a.R,
			"/",
			http.StatusFound,
		)
		return
	}

	hndl(a)
}}

func Greet(a muxes.HndlArg) {
	tmpls.Exec(a.W, "default", "greet", nil)
}

func main(){
	var err error

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
		FuncMap: templates.FuncMap{},
	}

	tmpls, err = templates.Parse(cfg)
	if err != nil {
		panic(err)
	}

	authorize := muxes.Chain{Authorize}
	unauthorize := muxes.Chain{Unauthorize}

	defs := []muxes.HndlDef {
		{
			"/", "^$",
			muxes.Handlers{
				"GET": muxes.Chained(authorize, Greet),
			},
		},
		{
			"/login/", "^$", muxes.Handlers {
				"GET": muxes.Chained(unauthorize, LoginGet),
				"POST": muxes.Chained(unauthorize, LoginPost),
			},
		},
		{"/get-test/", "", muxes.Handlers{"GET": muxes.GetTest} },
	}

	mux := muxes.Define(nil, defs)
	muxes.DefineStatic(mux, staticPath, "/s/")
	srv := http.Server {
		Addr: *AddrStr,
		Handler: mux,
	}

	usersJson, err := os.ReadFile(usersDbPath)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(usersJson, &users) ; err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", users)

	log.Printf("%s: Trying to run on '%s'...\n",
		os.Args[0],
		*AddrStr)
	log.Fatal(srv.ListenAndServe())
}
