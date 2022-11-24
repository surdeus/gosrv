package restx

import (
	"net/http"
	"github.com/surdeus/go-srv/src/dbx/sqlx"
	"github.com/surdeus/go-srv/src/httpx/muxx"
	//"github.com/surdeus/go-srv/src/urlx"
	"strings"
	"regexp"
	"fmt"
	"log"
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
		a muxx.HndlArg,
	) {
		tsName := a.R.URL.Path[len(pref):]
		tsName = strings.SplitN(tsName, "/", 2)[0]
		fmt.Println(tsName, mp)
		hndl, ok := mp[tsName]
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
) http.HandlerFunc {
	_, _, err := ts.PrimaryKeyFieldId()
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

func SqlParseParseValues(hndl muxx.Handler) muxx.Handler {
return func(a muxx.HndlArg) {
	hndl(a)
}}

func SqlMakeGetHandler(
	db *sqlx.DB,
	ts *sqlx.TableSchema,
	cfg *ArgCfg,
) muxx.Handler {
return func(a muxx.HndlArg) {
	fmt.Println("in getting handler")
	args := cfg.ParseValues(a.Values())
	fmt.Println(args)

	colArg, ok := args["c"]
	if !ok {
		a.NotFound()
	}

	columns := colArg.Values

	q := sqlx.Query{
		DB: db,
		Type: sqlx.SelectQueryType,
		Table: ts.Name,
		Columns: columns,
		/*Where: sqlx.Where {
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
		},*/
	}

	cs := sqlx.Conditions{}
	for k, arg := range args {
		if k == "c" {
			continue
		}

		name := arg.Splits[0]
		opStr := arg.Splits[1]

		op, _ := sqlx.
			ConditionOpStringMap[opStr]

		c := sqlx.Condition{
			Op: op,
			Values: [2]sqlx.RawValuer {
				sqlx.RawValue(name),
				sqlx.RawValue(arg.Values[0]),
			},
		}

		cs = append(cs, c)
	}

	q.Conditions = cs
	rows, err := q.Do()
	if err != nil {
		log.Println(err)
	} else {
		arr := make([]any, len(columns))
		for rows.Next() {
			rows.Scan(arr...)
			fmt.Println(arr)
		}
	}
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


