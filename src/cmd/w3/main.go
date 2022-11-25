package main

import(
	"os"
	"fmt"
	"flag"
	"log"
	"net/http"
	"encoding/json"
	"github.com/surdeus/go-srv/src/tmplx"
	"github.com/surdeus/go-srv/src/httpx/muxx"
	"github.com/surdeus/go-srv/src/httpx/cookiex"
	"github.com/surdeus/go-srv/src/authx"
	"github.com/surdeus/go-srv/src/httpx/restx"
	"github.com/surdeus/go-srv/src/dbx/sqlx"
	//"github.com/surdeus/go-srv/src/dbx/sqlx/qx"
	_ "github.com/go-sql-driver/mysql"
)
type Test struct {
	Id int
	Value int
	StringValue string
}

func (t Test)Sql() sqlx.TableSchema {
	return sqlx.TableSchema {
		Name: "Tests",
		Columns: sqlx.Columns {
			{
				Name: "Id",
				Type: "int(11)",
				Nullable: false,
				Key: "PRI",
				Extra: "auto_increment",
			},{
				OldName: "SuckValue",
				Name: "DickValue",
				Type: "int(11)",
				Nullable: true,
				Default: "25",
			},{
				Name: "StringValue",
				Type: "varchar(64)",
				Nullable: true,
				Default: "'some русская string'",
			},{
				Name: "NewValue",
				Type: "bigint(11)",
				Nullable: true,
				Default: "0",
			},
		},
	}
}

type AnotherTest struct {
	Id int
	Value int
	StringValue string
}

func (t AnotherTest)Sql() sqlx.TableSchema {
	return sqlx.TableSchema {
		OldName: "BetterTests",
		Name: "AnotherTests",
		Columns: sqlx.Columns {
			{
				Name: "Id",
				Type: "int(11)",
				Nullable: false,
				Key: "PRI",
				//Extra: "auto_increment",
			},{
				Name: "AnotherValue",
				Type: "int(11)",
				Nullable: true,
				Default: "25",
			},
		},
	}
}
type Token string
type Session struct {
	Reloaded int
	Email string
}
type Users map[string] string

var (
	sessions authx.Sessions
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

func HelloWorld(a muxx.HndlArg) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := "default"
	_, ok := a.Q["tmpl"]
	if ok {
		tmpl = a.Q["tmpl"][0]
	}
	tmpls.Exec(a.W, tmpl, "hellos/en", struct{}{})
}

func SalutonMondo(a muxx.HndlArg) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	name := "Mondo"
	_, ok := a.Q["name"]
	if ok {
		name = a.Q["name"][0]
	}
	tmpls.Exec(a.W, "default", "hellos/eo", struct{Name string}{Name: name})
}

func GetCookies(a muxx.HndlArg) {
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

func LoginGet(a muxx.HndlArg) {
	a.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpls.Exec(a.W, "unauth", "login", nil)
}

func LoginPost(a muxx.HndlArg) {
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

	token := sessions.New(formEmail)

	cookie := &http.Cookie{
		Name: "auth-token",
		Value: sessions.EncodeForClient(token),
		Path: "/",
	}

	http.SetCookie(a.W, cookie)
	http.Redirect(a.W, a.R, "/", http.StatusFound)
}

func Authorize(hndl muxx.Handler) muxx.Handler {
return func(a muxx.HndlArg) {
	cookie := a.R.Cookies()

	// No needed cookie, make user authorize.
	authToken, ok := cookiex.ByName(cookie, "auth-token")
	if !ok {
		http.Redirect(a.W, a.R,
			"/login/",
			http.StatusFound,
		)
		return
	}

	token, err := sessions.DecodeForServer(authToken)
	if err != nil {
		http.NotFound(a.W, a.R)
	}

	// No such token in sessions. Remove cookie and make authorize.
	email, loggedIn := sessions.Get(token)
	if !loggedIn {
		cookiex.Delete(a.W, "auth-token")
		http.Redirect(
			a.W,
			a.R,
			"/login/",
			http.StatusFound,
		)
		return
	}

	a.V["email"] = email
	hndl(a)
}}

func Unauthorize(hndl muxx.Handler) muxx.Handler {
return func(a muxx.HndlArg) {
	cookie := a.R.Cookies()

	_, ok := cookiex.ByName(cookie, "auth-token")
	if ok {
		http.Redirect(a.W, a.R,
			"/",
			http.StatusFound,
		)
		return
	}

	hndl(a)
}}

func Greet(a muxx.HndlArg) {
	email, _ := a.V["email"].(string)
	tmpls.Exec(a.W, "default", "greet",
		struct{
			Email string
		}{
			email,
		},
	)
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

	authorize := muxx.Chain{Authorize}
	unauthorize := muxx.Chain{Unauthorize}

	db, err := sqlx.Open(
		sqlx.ConnConfig{
			Driver: "mysql",
			Login: "test",
			Password: "hello",
			Host: "localhost",
			Port: 3306,
			Name: "test",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlers := []sqlx.Sqler{
		Test{},
		AnotherTest{},
	}
	db.Migrate(sqlers)

	defs := []muxx.HndlDef {
		{
			"/", "^$",
			muxx.Handlers{
				"GET": muxx.Chained(authorize, Greet),
			},
		},
		{
			"/login/", "^$", muxx.Handlers {
				"GET": muxx.Chained(unauthorize, LoginGet),
				"POST": muxx.Chained(unauthorize, LoginPost),
			},
		},
		{"/get-test/", "", muxx.Handlers{"GET": muxx.GetTest} },
		restx.Sql(db, "/api/", sqlers),
	}

	mux := muxx.Define(nil, defs)
	muxx.DefineStatic(mux, staticPath, "/s/")
	/*muxx.DefineSimple(
		mux,
		"/api/",
		restx.Sql(
			db,
			"/api/",
			sqlers,
		),
	)*/
	muxx.DefineSimple(
		mux,
		"/someshit/",
		func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "%s", "It works!")
		},
	)
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

	sessions = authx.New()

	fmt.Printf("%v\n", users)

	q := sqlx.Query{
		DB: db,
		Type: sqlx.SelectQueryType,
		Table: "Tests",
		Columns: []string{
			"DickValue",
			"StringValue",
		},
		Conditions: []sqlx.Condition {
			{
				Op: sqlx.GtConditionOp,
				Values: [2]sqlx.RawValuer{
					sqlx.RawValue("DickValue"),
					sqlx.Int(5),
				},
			},
			{
				Op: sqlx.EqConditionOp,
				Values: [2]sqlx.RawValuer{
					sqlx.RawValue("StringValue"),
					sqlx.String("value"),
				},
			},
		},
	}
	qs, err := q.Code()
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("%q\n", qs)
	}
	rows, err := q.Do()
	if err != nil {
		log.Println(err)
	} else {
		var (
			i int
			s string
		)
		fmt.Println("in")
		for rows.Next() {
			rows.Scan(&i, &s)
			fmt.Printf("%d %q\n", i, s)
		}
	}

	log.Printf("%s: Trying to run on '%s'...\n",
		os.Args[0],
		*AddrStr)
	log.Fatal(srv.ListenAndServe())
}
