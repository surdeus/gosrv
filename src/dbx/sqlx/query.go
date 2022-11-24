package sqlx

import (
	"fmt"
	"strings"
	"errors"
	"database/sql"
	"strconv"
)

// The interface type must implement to be converted to
// SQL code to be inserted.
type RawValuer interface {
	SqlRawValue() (RawValue, error)
}

type String string
type Int int
type Float32 float32
type Float64 float32
type Variable string

// Type to save values for substitution.
type RawValue string
// Type to save SQL-code.
type Code string

type ConditionOp int
type QueryType int


type Condition struct {
	Op ConditionOp
	Values [2]RawValuer
}
type Conditions []Condition

type Query struct {
	DB *DB
	Type QueryType
	Table string
	Columns []string
	Conditions Conditions
}

const (
	EqConditionOp ConditionOp = iota
	GtConditionOp
	LtConditionOp
	GeConditionOp
	LeConditionOp
	NeConditionOp
)

const (
	SelectQueryType QueryType = iota
	InsertQueryType
	DeleteQueryType
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

	// For the restx package.
	ConditionOpStringMap = map[string] ConditionOp {
		"eq" : EqConditionOp,
		"ne" : NeConditionOp,
		"gt" : GtConditionOp,
		"ge" : GeConditionOp,
		"lt" : LtConditionOp,
		"le" : LeConditionOp,
	}
)

func (w Conditions)Code() (Code, error) {
	if len(w) == 0 {
		return "", nil
	}

	ret := " where"
	for i, c := range w {
		op, ok := ConditionOpMap[c.Op]
		if !ok {
			return "", UnknownConditionOpErr
		}

		val1, err := c.Values[0].SqlRawValue()
		if err != nil {
			return "", err
		}
		val2, err := c.Values[1].SqlRawValue()
		if err != nil {
			return "", err
		}
		ret += fmt.Sprintf(
			" %s %s %s",
			val1,
			op,
			val2,
		)
		if i < len(w)-1 {
			ret += " and"
		}
	}
	return Code(ret), nil
}

func (q Query)Code() (String, error) {
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

	where, err := q.Conditions.Code()
	if err != nil {
		return String(""), err
	}

	switch q.Type {
	case SelectQueryType :
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

func (q Query)Do() (*sql.Rows, error) {
	qs, err := q.Code()
	if err != nil {
		return nil, err
	}

	return q.DB.Query(string(qs))
}

func (v RawValue)SqlRawValue() (RawValue, error) {
	return v, nil
}

func (i Int)SqlRawValue() (RawValue, error) {
	return RawValue(strconv.Itoa(int(i))), nil
}

func (s String)SqlRawValue() (RawValue, error) {
	ret := fmt.Sprintf("'%s'", s)
	return RawValue(ret), nil
}

