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
	SqlRawValue(db *DB) (RawValue, error)
}

// Any type that can be converted to code
type Coder interface {
	SqlCode(db *DB) (Code, error)
}

type RawValuers []RawValuer

type String string
type Int int
type Ints []Int
type Float32 float32
type Float64 float64
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
	Schema *TableSchema
	Columns ColumnNames
	Where Conditions
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
	NoQueryType QueryType = iota
	SelectQueryType
	InsertQueryType
	DeleteQueryType
	RenameTableQueryType
	CreateTableQueryType
	ModifyQueryType
)

var (
	NoTablesSpecifiedErr = errors.New("no table specified")
	NoColumnsSpecifiedErr = errors.New("no columns specified")
	UnknownQueryTypeErr = errors.New("unknown query type")
	UnknownConditionOpErr = errors.New("unknown condition operator")
	NoDBSpecifiedErr = errors.New("no database specified")
	NoSchemaSpecifiedErr = errors.New("no schema specified")
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

func Q() Query {
	return Query{}
}

func (db *DB)Q() Query {
	return Query{}.WithDB(db)
}

func (w Conditions)SqlCode(db *DB) (Code, error) {
	if len(w) == 0 {
		return "", nil
	
	}

	ret := " where"
	for i, c := range w {
		op, ok := ConditionOpMap[c.Op]
		if !ok {
			return "", UnknownConditionOpErr
		}

		val1, err := c.Values[0].SqlRawValue(db)
		if err != nil {
			return "", err
		}
		val2, err := c.Values[1].SqlRawValue(db)
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

func (q Query)RenameTable() Query {
	q.Type = RenameTableQueryType
	return q
}

func (q Query)SqlCode(db *DB) (Code, error) {
	var (
		ret string
		err error
	)
	switch q.Type {
	case SelectQueryType :
		if q.From == TableName("") {
			return "", NoTablesSpecifiedErr
		}

		columns, err := q.Columns.SqlRawValue(q.DB)
		if err != nil {
			return Code(""), err
		}

		from, err := q.From.SqlRawValue(q.DB)
		if err != nil {
			return Code(""), err
		}

		where, err := q.Where.SqlCode(q.DB)
		if err != nil {
			return Code(""), err
		}

		ret = fmt.Sprintf(
			"select %s from %s%s ;",
			columns,
			from,
			where,
		)
	case RenameTableQueryType :
		if q.From == "" || q.To == "" {
			return "", NoTablesSpecifiedErr
		}
		ret = fmt.Sprintf(
			"alter table %s rename %s ;",
			q.From,
			q.To,
		)
	case CreateTableQueryType :
		if q.Schema == nil {
			return "", NoSchemaSpecifiedErr
		}

		ret, err = db.
			TableCreationStringForSchema(q.Schema)
		if err != nil {
			return "", err
		}
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

func (q Query)WithSchema(schema *TableSchema) Query {
	q.Schema = schema
	return q
}

func (q Query)WithTo(to TableName) Query {
	q.To = to
	return q
}

func (q Query)WithWhere(where Conditions) Query {
	q.Where = where
	return q
}

func (q Query)Select() Query {
	return q.WithType(SelectQueryType)
}

func (q Query)Insert() Query {
	return q.WithType(InsertQueryType)
}

func (q Query)CreateTable() Query {
	return q.WithType(CreateTableQueryType)
}

func (q Query)Do() (*sql.Rows, error) {
	if q.DB == nil {
		return nil, NoDBSpecifiedErr
	}
	qs, err := q.SqlCode(q.DB)
	if err != nil {
		return nil, err
	}

	return q.DB.Query(string(qs))
}

func (v TableName)SqlRawValue(db *DB) (RawValue, error) {
	return RawValue(v), nil
}

func (v ColumnName)SqlRawValue(db *DB) (RawValue, error) {
	return RawValue(v), nil
}

func (v RawValue)SqlRawValue(db *DB) (RawValue, error) {
	return v, nil
}

func (i Int)SqlRawValue(db *DB) (RawValue, error) {
	return RawValue(strconv.Itoa(int(i))), nil
}

func (tn TableNames)SqlRawValue(db *DB) (RawValue, error) {
	if len(tn) == 0 {
		return RawValue(""), NoTablesSpecifiedErr
	}

	buf := make([]string, 0)
	for _, t := range tn {
		v, err := t.SqlRawValue(db)
		if err != nil {
			return RawValue(""), err
		}
		buf = append(buf, string(v))
	}

	ret := strings.Join(buf, ", ")
	return RawValue(ret), nil
}

func (cn ColumnNames)SqlRawValue(db *DB) (RawValue, error) {
	if len(cn) == 0 {
		return RawValue(""), NoColumnsSpecifiedErr
	}

	buf := make([]string, 0)
	for _, c := range cn {
		v, err := c.SqlRawValue(db)
		if err != nil {
			return RawValue(""), err
		}
		buf = append(buf, string(v))
	}

	ret := strings.Join(buf, ", ")
	return RawValue(ret), nil
}

func (s String)SqlRawValue(db *DB) (RawValue, error) {
	ret := strings.ReplaceAll(string(s), "'", "''")
	ret = fmt.Sprintf("'%s'", ret)
	return RawValue(ret), nil
}

func (rvs RawValuers) SqlMultiValue(db *DB) (Code, error) {
	var ret Code
	for i, v := range rvs {
		raw, err := v.SqlRawValue(db)
		if err != nil {
			return "", err
		}

		ret += Code(raw)

		if i != len(rvs) - 1 {
			ret += ","
		}
	}

	return ret, nil
}

func (rvs RawValuers) SqlCodeTuple(db *DB) (Code, error) {
	mval, err := rvs.SqlMultiValue(db)
	if err != nil {
		return Code(""), err
	}

	if mval == "" {
		return "", nil
	}

	return Code(fmt.Sprintf("(%s)", mval)), nil
}

func (db *DB)RawValuersEq(
	v1, v2 RawValuer,
) (bool, error) {

	if v1 == nil && v2 == nil {
		return true, nil
	}

	if v1 == nil || v2 == nil {
		fmt.Println("in")
		return false, nil
	}

	raw1, err := v1.SqlRawValue(db)
	if err != nil {
		return false, err
	}

	raw2, err := v2.SqlRawValue(db)
	if err != nil {
		return false, err
	}
	return raw1 == raw2, nil
}

func (c Code)SqlCode(db *DB) (Code, error) {
	return c, nil
}

func (db *DB)CodersEq(
	c1, c2 Coder,
) (bool, error) {
	if c1 == nil && c2 == nil {
		return true, nil
	}

	if c1 == nil || c2 == nil {
		return false, nil
	}

	code1, err := c1.SqlCode(db)
	if err != nil {
		return false, err
	}

	code2, err := c2.SqlCode(db)
	if err != nil {
		return false, err
	}

	return code1 == code2, nil
}

