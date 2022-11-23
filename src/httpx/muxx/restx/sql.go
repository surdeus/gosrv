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
		"GET" : SqlMakeGetHandler(ts),
		"POST" : SqlMakePostHandler(ts),
		"PUT" : SqlMakePutHandler(ts),
		"PATCH" : SqlMakePatchHandler(ts),
		"DELETE" : SqlMakeDeleteHandler(ts),
	}
	return muxx.MakeHttpHandleFunc(
		pref,
		regexp.MustCompile("^[0-9]*$"),
		handlers,
	)
}

func SqlMakeGetHandler(
	ts *sqlx.TableSchema,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in getting handler")
}}

func SqlMakePostHandler(
	ts *sqlx.TableSchema,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in posting handler")
}}

func SqlMakePutHandler(
	ts *sqlx.TableSchema,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in putting handler")
}}

func SqlMakePatchHandler(
	ts *sqlx.TableSchema,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in patching handler")
}}

func SqlMakeDeleteHandler(
	ts *sqlx.TableSchema,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in deleting handler")
}}


