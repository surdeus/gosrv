package restx

import (
	"strings"
	"net/url"
	"errors"
	"github.com/surdeus/go-srv/src/dbx/sqlx"
	"fmt"
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

func (args Args)SqlGetQuery(
	ts *sqlx.TableSchema,
	tsMap map[sqlx.ColumnName] *sqlx.Column,
) (sqlx.Query, error) {
	columns, err := args.SqlColumns(ts, tsMap)
	if err != nil {
		return sqlx.Q(), err
	}

	cs, err := args.SqlConditions()
	if err != nil {
		return sqlx.Q(), err
	}

	fmt.Println("reached")
	fmt.Println(columns)
	q := sqlx.Query{
		From: ts.Name,
		ColumnNames: columns,
		Where: cs,
	}

	return q, nil
}

func (args Args)SqlColumns(
	ts *sqlx.TableSchema,
	tsMap map[sqlx.ColumnName] *sqlx.Column,
) (sqlx.ColumnNames, error) {
	colArg, ok := args["c"]
	if !ok {
		return sqlx.ColumnNames{}, NoColumnsSpecifiedErr
	}

	columnsStr := colArg.Values
	columns := sqlx.ColumnNames{}
	for _, c := range columnsStr {
		if c == "*" {
			return ts.Columns.Names(), nil
		}
		_, exists := tsMap[sqlx.ColumnName(c)]
		if !exists {
			return sqlx.ColumnNames{},
				sqlx.ColumnDoesNotExistErr
		}
		columns = append(
			columns,
			sqlx.ColumnName(c),
		)
	}

	return columns, nil
}

func (args Args)SqlConditions() (sqlx.Conditions, error) {
	cs := sqlx.Conditions{}
	for k, arg := range args {
		if k == "c" {
			continue
		}

		if len(arg.Splits) != 2 {
			return sqlx.Conditions{},
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

	return cs, nil
}
