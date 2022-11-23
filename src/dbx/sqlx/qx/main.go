package qx

import (
	"github.com/surdeus/go-srv/src/dbx/sqlx"
	"fmt"
	"strings"
	"errors"
)

type String string
type ConditionOp int
type Type int
type Value string

type Condition struct {
	Op ConditionOp
	Values [2]Value
}

type Where struct {
	Conditions []Condition
}

type Query struct {
	db *sqlx.DB
	Type Type
	Table string
	Columns []string
	Where Where
}

const (
	EqConditionOp ConditionOp = iota
	GtConditionOp
	LtConditionOp
	GeConditionOp
	LeConditionOp
)

const (
	SelectType Type = iota
	InsertType
	DeleteType
)

var (
	NoTableSpecifiedErr = errors.New("no table specified")
	NoColumnsSpecifiedErr = errors.New("no columns specified")
	UnknownQueryTypeErr = errors.New("unknown query type")
	UnknownConditionOpErr = errors.New("unknown condition operator")

	ConditionOpMap = map[ConditionOp] String {
		EqConditionOp : "=",
		NeConditionOp : "<>",
		GtConditionOp : ">",
		GeConditionOp : ">=",
		LtConditionOp : "<",
		LeConditionOp : ">=",
	}
)

func (w Where)SqlString() (String, error) {
	if len(w.Conditions) == 0 {
		return String(""), nil
	}

	ret := " where"
	for i, c := range w.Conditions {
		op, ok := ConditionOpMap[c.Op]
		if !ok {
			return "", UnknownConditionOpErr
		}

		ret += fmt.Sprintf(
			" %s %s %s",
			c.Values[0],
			op,
			c.Values[1],
		)
		if i < len(w.Conditions)-1 {
			ret += " and"
		}
	}
	return String(ret), nil
}

func (q Query)SqlString() (String, error) {
	var (
		ret, c string
	)

	if q.Table == "" {
		return "", NoTableSpecifiedErr
	}

	if len(q.Columns) == 0 {
		return "", NoColumnsSpecifiedErr
	}

	c = strings.Join(q.Columns, ", ")

	where, err := q.Where.SqlString()
	if err != nil {
		return String(""), err
	}

	switch q.Type {
	case SelectType :
		ret = fmt.Sprintf(
			"select %s from %s%s;",
			c,
			q.Table,
			where,
		)
	default:
		return "", UnknownQueryTypeErr
	}

	return String(ret), nil
}

