package restx

import (
	"net/http"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	"github.com/surdeus/gosrv/src/httpx/muxx"
	//"github.com/surdeus/godat/src/slicex"
	//"github.com/surdeus/gosrv/src/urlx"
	"strings"
	"regexp"
	"fmt"
	"encoding/json"
	//"reflect"
)

func Sql(
	db *sqlx.Db,
	pref string,
) muxx.HndlDef {
	ret := muxx.HndlDef{}
	ret.Pref = pref
	ret.Re = ""
	ret.Handlers = muxx.Handlers {
		"" : SqlHandler(db, pref),
	}

	return ret
}

// SQL REST access implementetaion
// based on sqlx package migration description.
func SqlHandler(
	db *sqlx.Db,
	pref string,
) muxx.Handler {
	schemas := db.Tables

	mp := make(map[sqlx.TableName] http.HandlerFunc)
	for _, schema := range schemas {
		mp[schema.Name] = MakeSqlTableHandler(
			db,
			pref + string(schema.Name) + "/",
			schema,
		)
	}
	return func(
		a muxx.HndlArg,
	) {
		tsName := a.R.URL.Path[len(pref):]
		tsName = strings.SplitN(tsName, "/", 2)[0]
		hndl, ok := mp[sqlx.TableName(tsName)]
		if !ok {
			a.NotFound()
			return
		}
		hndl(a.W, a.R)
	}
}

func MakeSqlTableHandler(
	db *sqlx.Db,
	pref string,
	ts *sqlx.TableSchema,
) http.HandlerFunc {
	_, _, err := ts.PrimaryKeyColumn()
	if err != nil {
		panic(err)
	}

	cfg := StdArgCfg()
	handlers := muxx.Handlers {
		"GET" : SqlMakeGetHandler(db, ts, cfg),
		"POST" : SqlMakePostHandler(db, ts, cfg),
		"PUT" : SqlMakePutHandler(db, ts, cfg),
		"PATCH" : SqlMakePatchHandler(db, ts, cfg),
		"DELETE" : SqlMakeDeleteHandler(db, ts, cfg),
	}

	fin := muxx.MakeHttpHandleFunc(
		pref,
		regexp.MustCompile("^[0-9]*$"),
		handlers,
	)

	return fin
}

func SqlMakeGetHandler(
	db *sqlx.Db,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
) muxx.Handler {
return func(a muxx.HndlArg) {
	cMap := db.TCMap[ts.Name]
	args := cfg.ParseValues(a.Values())
	
	q, err := args.SqlGetQuery(ts, cMap)
	if err != nil {
		a.ServerError(err)
		return
	}

	_, rows, err := db.Do(q)
	if err != nil {
		a.ServerError(err)
		return
	}
	defer rows.Close()

	if err != nil {
		a.ServerError(err)
		return
	}

	values, err := db.ReadRowValues(
		rows,
		q.GetTableName(),
		q.GetColumnNames(),
	)
	if err != nil {
		a.ServerError(err)
		return
	}


	ret := []any{}
	for v := range values {
		ret = append(ret, v)
	}

	bts, err := json.Marshal(ret)
	if err != nil {
		a.NotFound()
		return
	}

	a.W.Header().Set(
		"Content-Type",
		"application/json",
	)
	a.W.Write(bts)
}}

func SqlMakePostHandler(
	db *sqlx.Db,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in posting handler")
}}

func SqlMakePutHandler(
	db *sqlx.Db,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in putting handler")
}}

func SqlMakePatchHandler(
	db *sqlx.Db,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in patching handler")
}}

func SqlMakeDeleteHandler(
	db *sqlx.Db,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in deleting handler")
}}


