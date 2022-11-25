package restx

import (
	"strings"
	"net/url"
	"errors"
	"github.com/surdeus/go-srv/src/dbx/sqlx"
	//"fmt"
)

type ArgCfg struct {
	Sep string
}

type Args map[string] Arg

type Arg struct {
	Splits []string
	Values []string
}

var (
	NoColumnsSpecifiedErr = errors.New(
		"no columns specified",
	)
	WrongSplitOperatorFormatErr = errors.New(
		"wrong split operator format",
	)
)

func StdArgCfg() *ArgCfg {
	ret := ArgCfg {
		Sep : "__",
	}
	return &ret
}

func (ac *ArgCfg) ParseValues (
	values url.Values,
) Args {
	ret := make(Args)
	for k, v := range values {
		//fmt.Printf("%s %q, %q\n", k, v, strings.Split(k, ac.Sep))
		buf := Arg{}
		buf.Splits = strings.Split(k, ac.Sep)
		buf.Values = v
		ret[buf.Splits[0]] = buf
	}

	return ret
}

func (args Args)SqlQuery(ts *sqlx.TableSchema) (sqlx.Query, error) {
	colArg, ok := args["c"]
	if !ok {
		return sqlx.Query{}, NoColumnsSpecifiedErr
	}

	columns := colArg.Values

	q := sqlx.Query{
		Table: ts.Name,
		Columns: columns,
	}

	cs := sqlx.Conditions{}
	for k, arg := range args {
		if k == "c" {
			continue
		}

		if len(arg.Splits) != 2 {
			return sqlx.Query{},
				WrongSplitOperatorFormatErr 
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

	return q, nil
}
