package restx

import (
	"net/http"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	"github.com/surdeus/gosrv/src/httpx/muxx"
	"github.com/surdeus/godat/src/slicex"
	//"github.com/surdeus/gosrv/src/urlx"
	"strings"
	"regexp"
	"fmt"
	"log"
	"encoding/json"
	//"reflect"
)

func Sql(
	db *sqlx.Db,
	pref string,
	sqlers []any,
) muxx.HndlDef {
	ret := muxx.HndlDef{}
	ret.Pref = pref
	ret.Re = ""
	ret.Handlers = muxx.Handlers {
		"" : SqlHandler(db, pref, sqlers),
	}

	return ret
}

// SQL REST access implementetaion
// based on sqlx package migration description.
func SqlHandler(
	db *sqlx.Db,
	pref string,
	rcs []any,
) muxx.Handler {
	schemas := sqlx.TableSchemas{}
	for _, sqler := range rcs {
		ts := sqler.(sqlx.Sqler).Sql()
		schemas = append(schemas, ts)
	}

	mp := make(map[sqlx.TableName] http.HandlerFunc)
	for i, schema := range schemas {
		mp[schema.Name] = MakeSqlTableHandler(
			db,
			pref + string(schema.Name) + "/",
			schema,
			rcs[i],
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
	rc any,
) http.HandlerFunc {
	_, _, err := ts.PrimaryKeyColumn()
	if err != nil {
		panic(err)
	}

	tsMap := slicex.MakeMap(
		ts.Columns,
		func(cols []*sqlx.Column, i int) sqlx.ColumnName {
			return cols[i].Name
		},
	)
	cfg := StdArgCfg()
	handlers := muxx.Handlers {
		"GET" : SqlMakeGetHandler(db, ts, cfg, tsMap, rc),
		"POST" : SqlMakePostHandler(db, ts, cfg, rc),
		"PUT" : SqlMakePutHandler(db, ts, cfg, rc),
		"PATCH" : SqlMakePatchHandler(db, ts, cfg, rc),
		"DELETE" : SqlMakeDeleteHandler(db, ts, cfg, rc),
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
	tsMap map[sqlx.ColumnName] *sqlx.Column,
	rc any,
) muxx.Handler {
return func(a muxx.HndlArg) {
	args := cfg.ParseValues(a.Values())
	
	q, err := args.SqlGetQuery(ts, tsMap)
	if err != nil {
		log.Println(err)
		a.NotFound()
		return
	}

	_, rows, err := db.Do(q)
	if err != nil {
		log.Println(err)
		a.NotFound()
		return
	}
	defer rows.Close()

	if err != nil {
		a.NotFound()
		return
	}

	values, err := db.ReadRowValues(
		rows,
		ts,
		q.GetColumns(),
		tsMap,
		rc,
	)
	if err != nil {
		a.NotFound()
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
	rc any,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in posting handler")
}}

func SqlMakePutHandler(
	db *sqlx.Db,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
	rc any,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in putting handler")
}}

func SqlMakePatchHandler(
	db *sqlx.Db,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
	rc any,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in patching handler")
}}

func SqlMakeDeleteHandler(
	db *sqlx.Db,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
	rc any,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in deleting handler")
}}


