package restx

import (
	"net/http"
	"github.com/surdeus/go-srv/src/dbx/sqlx"
	"github.com/surdeus/go-srv/src/httpx/muxx"
	"github.com/surdeus/godat/src/slicex"
	//"github.com/surdeus/go-srv/src/urlx"
	"strings"
	"regexp"
	"fmt"
	"log"
	"encoding/json"
	//"reflect"
)

func Sql(
	db *sqlx.DB,
	pref string,
	sqlers []sqlx.Sqler,
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
	sqlers []sqlx.Sqler,
) muxx.Handler {
	schemas := sqlx.TableSchemas{}
	for _, sqler := range sqlers {
		ts := sqler.Sql()
		schemas = append(schemas, ts)
	}

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
		fmt.Println(tsName, mp)
		hndl, ok := mp[sqlx.TableName(tsName)]
		if !ok {
			fmt.Println("in")
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
	fmt.Println(tsMap)
	cfg := StdArgCfg()
	handlers := muxx.Handlers {
		"GET" : SqlMakeGetHandler(db, ts, cfg, tsMap),
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

func SqlParseParseValues(hndl muxx.Handler) muxx.Handler {
return func(a muxx.HndlArg) {
	hndl(a)
}}

func SqlMakeGetHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
	tsMap map[sqlx.ColumnName] *sqlx.Column,
) muxx.Handler {
return func(a muxx.HndlArg) {
	//fmt.Println("in getting handler")
	args := cfg.ParseValues(a.Values())
	//fmt.Println(args)
	
	q, err := args.SqlGetQuery(ts)
	if err != nil {
		log.Println(err)
		a.NotFound()
		return
	}
	q = q.WithDB(db).
		WithType(sqlx.SelectQueryType)

	s, err := q.SqlCode(db)
	if err != nil {
		log.Println(err)
		a.NotFound()
		return
	}
	println(s)

	rows, err := q.Do()
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	ret := []any{}
	row := []any{}
	for _, v := range q.ColumnNames {
		c, ok := tsMap[v]
		if !ok {
			a.NotFound()
			return
		}
		switch c.Type.VarType {
		case sqlx.VarcharColumnVarType :
			row = append(row, new(string))
		case sqlx.IntColumnVarType :
			row = append(row, new(int))
		default:
			a.NotFound()
			return
		}
	}
	for rows.Next() {
		err = rows.Scan(row...)
		fmt.Println(err)
		add := []any{}
		for _, v := range row {
			switch v.(type) {
			case *int :
				add = append(add, *(v.(*int)))
			case *string :
				println("in")
				add = append(add, *(v.(*string)))
			}
		}
		ret = append(ret, add)
	}

	fmt.Printf("%q\n", ret)
	js, err := json.Marshal(ret)
	if err != nil {
		http.Error(a.W, err.Error(), http.StatusInternalServerError)
		return
	}

	a.W.Header().Set(
		"Content-Type",
		"application/json",
	)
	a.W.Write(js)
}}

func SqlMakePostHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in posting handler")
}}

func SqlMakePutHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in putting handler")
}}

func SqlMakePatchHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in patching handler")
}}

func SqlMakeDeleteHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in deleting handler")
}}


