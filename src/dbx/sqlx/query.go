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

type ColumnName string
type ColumnNames []ColumnName
type TableName string 
type TableNames []TableName

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
	From TableName
	To TableName
	Columns ColumnNames
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
	AlterTableRenameQueryType
	ModifyQueryType
)

var (
	NoTablesSpecifiedErr = errors.New("no table specified")
	NoColumnsSpecifiedErr = errors.New("no columns specified")
	UnknownQueryTypeErr = errors.New("unknown query type")
	UnknownConditionOpErr = errors.New("unknown condition operator")

	ConditionOpMap = map[ConditionOp] String {
		EqConditionOp : "=",
		NeConditionOp : "<>",
		GtConditionOp : ">",
		GeConditionOp : ">=",
		LtConditionOp : "<",
		LeConditionOp : "<=",
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

func (q Query)Code() (Code, error) {
	var (
		ret string
	)
	switch q.Type {
	case SelectQueryType :
		if q.From == TableName("") {
			return "", NoTablesSpecifiedErr
		}

		columns, err := q.Columns.SqlRawValue()
		if err != nil {
			return Code(""), err
		}

		from, err := q.From.SqlRawValue()
		if err != nil {
			return Code(""), err
		}

		where, err := q.Conditions.Code()
		if err != nil {
			return Code(""), err
		}

		ret = fmt.Sprintf(
			"select %s from %s%s ;",
			columns,
			from,
			where,
		)
	case AlterTableRenameQueryType :
		ret = fmt.Sprintf(
			"alter table %s rename %s ;",
			q.From,
			q.To,
		)
	default:
		return "", UnknownQueryTypeErr
	}

	return Code(ret), nil
}

func (q Query)WithDB(db *DB) Query {
	q.DB = db
	return q
}

func (q Query)WithType(t QueryType) Query {
	q.Type = t
	return q
}

func (q Query)WithFrom(from TableName) Query {
	q.From = from
	return q
}

func (q Query)WithTo(to TableName) Query {
	q.To = to
	return q
}

func (q Query)AlterTableRename() Query {
	return q.WithType(AlterTableRenameQueryType)
}

func (q Query)Select() Query {
	return q.WithType(SelectQueryType)
}

func (q Query)Insert() Query {
	return q.WithType(InsertQueryType)
}

func (q Query)Do() (*sql.Rows, error) {
	qs, err := q.Code()
	if err != nil {
		return nil, err
	}

	return q.DB.Query(string(qs))
}

func (v TableName)SqlRawValue() (RawValue, error) {
	return RawValue(v), nil
}

func (v ColumnName)SqlRawValue() (RawValue, error) {
	return RawValue(v), nil
}

func (v RawValue)SqlRawValue() (RawValue, error) {
	return v, nil
}

func (i Int)SqlRawValue() (RawValue, error) {
	return RawValue(strconv.Itoa(int(i))), nil
}

func (tn TableNames)SqlRawValue() (RawValue, error) {
	if len(tn) == 0 {
		return RawValue(""), NoTablesSpecifiedErr
	}

	buf := make([]string, 0)
	for _, t := range tn {
		v, err := t.SqlRawValue()
		if err != nil {
			return RawValue(""), err
		}
		buf = append(buf, string(v))
	}

	ret := strings.Join(buf, ", ")
	return RawValue(ret), nil
}

func (cn ColumnNames)SqlRawValue() (RawValue, error) {
	if len(cn) == 0 {
		return RawValue(""), NoColumnsSpecifiedErr
	}

	buf := make([]string, 0)
	for _, c := range cn {
		v, err := c.SqlRawValue()
		if err != nil {
			return RawValue(""), err
		}
		buf = append(buf, string(v))
	}

	ret := strings.Join(buf, ", ")
	return RawValue(ret), nil
}

func (s String)SqlRawValue() (RawValue, error) {
	ret := strings.ReplaceAll(string(s), "'", "''")
	ret = fmt.Sprintf("'%s'", ret)
	return RawValue(ret), nil
}

