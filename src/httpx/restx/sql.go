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
	//"encoding/json"
	"encoding/gob"
	//"reflect"
)

func Sql(
	db *sqlx.DB,
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
	db *sqlx.DB,
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
	db *sqlx.DB,
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
	db *sqlx.DB,
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
	q = q.WithDB(db).
		WithType(sqlx.SelectQueryType)

	rows, err := q.Do()
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
		q.ColumnNames,
		tsMap,
		rc,
	)

	a.W.Header().Set(
		"Content-Type",
		"application/gob",
	)

	
	enc := gob.NewEncoder(a.W)
	for v := range values {
		enc.Encode(v)
	}
}}

func SqlMakePostHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
	rc any,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in posting handler")
}}

func SqlMakePutHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
	rc any,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in putting handler")
}}

func SqlMakePatchHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
	rc any,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in patching handler")
}}

func SqlMakeDeleteHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
	rc any,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in deleting handler")
}}


