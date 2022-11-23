package restx

import (
	"net/http"
	"github.com/surdeus/go-srv/src/dbx/sqlx"
	"github.com/surdeus/go-srv/src/httpx/muxx"
	"strings"
	"regexp"
	"fmt"
)

// SQL REST access implementetaion
// based on sqlx package migration description.
func Sql(
	db *sqlx.DB,
	pref string,
	sqlers []sqlx.Sqler,
) http.HandlerFunc {
	schemas := []*sqlx.TableSchema{}
	for _, sqler := range sqlers {
		ts := sqler.Sql()
		schemas = append(schemas, &ts)
	}

	mp := make(map[string] http.HandlerFunc)
	for _, schema := range schemas {
		mp[schema.Name] = MakeSqlTableHandler(
			db,
			pref + schema.Name + "/",
			schema,
		)
	}
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		tsName := r.URL.Path[len(pref):]
		tsName = strings.SplitN(tsName, "/", 2)[0]
		fmt.Println(tsName, mp)
		hndl, ok := mp[tsName]
		if !ok {
			http.NotFound(w, r)
			return
		}
		hndl(w, r)
	}
}

func MakeSqlTableHandler(
	db *sqlx.DB,
	pref string,
	ts *sqlx.TableSchema,
) http.HandlerFunc {
	_, _, err := ts.PrimaryKeyFieldId()
	if err != nil {
		panic(err)
	}
	handlers := muxx.Handlers {
		"GET" : SqlMakeGetHandler(db, ts),
		"POST" : SqlMakePostHandler(db, ts),
		"PUT" : SqlMakePutHandler(db, ts),
		"PATCH" : SqlMakePatchHandler(db, ts),
		"DELETE" : SqlMakeDeleteHandler(db, ts),
	}
	return muxx.MakeHttpHandleFunc(
		pref,
		regexp.MustCompile("^[0-9]*$"),
		handlers,
	)
}

func SqlMakeGetHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in getting handler")
	//mp := a.R.URL.Query()
}}

func SqlMakePostHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in posting handler")
}}

func SqlMakePutHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in putting handler")
}}

func SqlMakePatchHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in patching handler")
}}

func SqlMakeDeleteHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in deleting handler")
}}


