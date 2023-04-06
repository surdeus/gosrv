package main

import(
	"os"
	"fmt"
	"flag"
	"log"
	"net/http"
	"encoding/json"
	"github.com/surdeus/gosrv/src/tmplx"
	"github.com/surdeus/gosrv/src/httpx"
	"github.com/surdeus/gosrv/src/httpx/authx"
	//"github.com/surdeus/gosrv/src/httpx/restx"
	//"github.com/surdeus/gosrv/src/dbx/sqlx"
	//"github.com/surdeus/gosrv/src/httpx/apix"
	//dbtest "github.com/surdeus/gosrv/src/cmd/dbtest/structs"
	//"github.com/surdeus/gosrv/src/dbx/sqlx/qx"
	//_ "github.com/go-sql-driver/mysql"
)

type Token string
type Session struct {
	Reloaded int
	Email string
}
type Users map[string] string

var (
	sessions authx.Sessions[*Session]
	tokens = make(map[string] string)
	tmpls tmplx.Templates
	users Users
	datPath = "dat/"
	dbPath = datPath+"db/"
	staticPath = datPath+"s/"
	usersDbPath = dbPath + "users"
)

type ContextKey string
const ContextEmailKey ContextKey = "email"

func HelloWorld(a *httpx.Context) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := "default"
	_, ok := a.Q["tmpl"]
	if ok {
		tmpl = a.Q["tmpl"][0]
	}
	tmpls.Exec(a.W, tmpl, "hellos/en", struct{}{})
}

func SalutonMondo(a *httpx.Context) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	name := "Mondo"
	_, ok := a.Q["name"]
	if ok {
		name = a.Q["name"][0]
	}
	tmpls.Exec(a.W, "default", "hellos/eo", struct{Name string}{Name: name})
}

func GetCookies(a *httpx.Context) {
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

func LoginGet(a *httpx.Context) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpls.Exec(a.W, "unauth", "login", nil)
}

func LoginPost(a *httpx.Context) {
	formEmail := a.R.Form.Get("email")
	formPassword := a.R.Form.Get("password")
	fmt.Printf(
		"Email: '%s', password: '%s'\n",
		formEmail,
		formPassword,
	)

	password, ok := users[formEmail]
	if !ok {
		a.NotFound()
		return
	} else if password != formPassword  {
		a.W.WriteHeader(http.StatusUnauthorized)
		a.W.Write([]byte("Unauth"))
		return
	}

	token := sessions.New(&Session{
		Email: formEmail,
		Reloaded: 0,
	})

	cookie := &http.Cookie{
		Name: "auth-token",
		Value: sessions.EncodeForClient(token),
		Path: "/",
	}

	http.SetCookie(a.W, cookie)
	http.Redirect(a.W, a.R, "/", http.StatusFound)
}

func Authorize(hndl httpx.HandlerFunc) httpx.HandlerFunc {
return func(a *httpx.Context) {
	cookies := a.Cookies()

	// No needed cookie, make user authorize.
	authToken, ok := cookies.Get("auth-token")
	if !ok {
		a.Redirect( "/login/", http.StatusFound)
		return
	}

	token, err := sessions.DecodeForServer(authToken)
	if err != nil {
		http.NotFound(a.W, a.R)
	}

	// No such token in sessions. Remove cookie and make authorize.
	email, loggedIn := sessions.Get(token)
	if !loggedIn {
		a.DeleteCookie("auth-token")
		a.Redirect("/login/", http.StatusFound)
		return
	}

	a.V["email"] = email
	hndl(a)
}}

func Unauthorize(hndl httpx.HandlerFunc) httpx.HandlerFunc {
return func(a *httpx.Context) {
	cookies := a.Cookies()

	_, ok := cookies.Get("auth-token")
	if ok {
		a.Redirect("/", http.StatusFound)
		return
	}

	hndl(a)
}}

func Greet(a *httpx.Context) {
	session, _ := a.V["email"].(*Session)
	tmpls.Exec(a.W, "default", "greet", session)
	session.Reloaded++
}

func main(){
	var err error

	AddrStr := flag.String("a", ":8080", "Adress string")
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		os.Exit(1)
	}

	funcCfg := tmplx.StdFuncCfg()

	fmap := tmplx.StdFuncMap()
	fmap["styles"] = funcCfg.Styles
	fmap["scripts"] = funcCfg.Scripts
	cfg := tmplx.ParsingConfig{
		Component: "tmpl/c/",
		View: "tmpl/v/",
		Template: "tmpl/t/",
		FuncMap: fmap,
	}

	tmpls, err = tmplx.Parse(cfg)
	if err != nil {
		panic(err)
	}

	authorize := httpx.Chain{Authorize}
	unauthorize := httpx.Chain{Unauthorize}

	/*sqlers := []sqlx.Sqler{
		dbtest.Test{},
		dbtest.AnotherTest{},
	}
	db, err := sqlx.Open(
		sqlx.ConnConfig{
			Driver: "mysql",
			Login: "test",
			Password: "hello",
			Host: "localhost",
			Port: 3306,
			Name: "test",
		},
		sqlers,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Migrate()
	if err != nil {
		panic(err)
	}*/

	mux := httpx.NewMux()
	mux.Define(
		httpx.Def("/").Re("$^").
			Method("GET", httpx.Chained(authorize, Greet)),
		httpx.Def("/login/").Re("$^").
			Method("GET", httpx.Chained(unauthorize, LoginGet)).
			Method("POST", httpx.Chained(unauthorize, LoginPost)),
		httpx.Def("/get-test/").Re("").
			Method("GET", httpx.GetTest),
		httpx.Def("/s/").StaticFiles(staticPath),
		httpx.Def("/someshit/").SimpleHandlerFunc(
		func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "%s", "It works!")
		}),
	)
	/*httpx.DefineSimple(
		mux,
		"/api/",
		restx.Sql(
			db,
			"/api/",
			sqlers,
		),
	)*/
	srv := http.Server {
		Addr: *AddrStr,
		Handler: mux,
	}

	usersJson, err := os.ReadFile(usersDbPath)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(
		usersJson,
		&users,
	) ; err != nil {
		panic(err)
	}

	sessions = authx.New[*Session]()
	fmt.Printf("%v\n", users)

	log.Printf("%s: Trying to run on '%s'...\n",
		os.Args[0],
		*AddrStr)
	log.Fatal(srv.ListenAndServe())
}
