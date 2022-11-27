package sqlx

import (
	"fmt"
	"strings"
	"errors"
	"database/sql"
	"strconv"
)

// The interface type must implement to be converted to
// SQL code to be inserted for safety.
type Rawer interface {
	SqlRaw(db *DB) (Raw, error)
}

type Rawers []Rawer

type String string
type Int int
type Ints []Int
type Double float64
type Variable string
type Null struct{}

type ColumnName string
type ColumnNames []ColumnName
type TableName string 
type TableNames []TableName

// Type to save raw strings for substitution.
type Raw string

type ConditionOp int
type QueryType int

type Condition struct {
	Op ConditionOp
	Values [2]Rawer
}
type Conditions []Condition
type Where Conditions
type Wheres []Where

type Query struct {
	DB *DB
	Type QueryType
	TableSchemas TableSchemas
	ColumnNames ColumnNames
	TableNames TableNames
	Columns Columns
	ColumnTypes []ColumnType
	Wheres Wheres
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
	RenameColumnQueryType
	CreateTableQueryType
	AlterColumnTypeQueryType
	ModifyQueryType
)

var (
	NoTablesSpecifiedErr = errors.New("no table specified")
	NoColumnsSpecifiedErr = errors.New("no columns specified")
	WrongNumOfColumnsSpecifiedErr = errors.New(
		"wrong number of columns specified")
	WrongQueryInputFormatErr = errors.New(
		"wrong query input format",
	)
	WrongValuerFormatErr = errors.New("wrong valuer format")
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

func (w Where)SqlRaw(db *DB) (Raw, error) {
	if len(w) == 0 {
		return "", nil
	}
	ret, err := db.Rprintf(
		"where %s",
		Conditions(w),
	)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (w Conditions)SqlRaw(db *DB) (Raw, error) {
	if len(w) == 0 {
		return "", nil
	}

	ret := Raw("")
	for i, c := range w {
		op, ok := ConditionOpMap[c.Op]
		if !ok {
			return "", UnknownConditionOpErr
		}

		cond, err := db.Rprintf(
			"%s%s %s %s",
			Raw(" "),
			c.Values[0],
			op,
			c.Values[1],
		)
		if err != nil {
			return "", err
		}
		ret += cond

		if i < len(w)-1 {
			ret += Raw(" and")
		}
	}
	return Raw(ret), nil
}

func (q Query)SqlRaw(db *DB) (Raw, error) {
	var (
		ret Raw
		err error
	)
	if db != nil {
		q.DB = db
	}
	switch q.Type {
	case SelectQueryType :
		if len(q.TableNames) != 1 {
			return "", NoTablesSpecifiedErr
		}

		if len(q.ColumnNames) == 0 {
			return "", NoColumnsSpecifiedErr
		}

		if len(q.Wheres) > 1 {
			return "", WrongQueryInputFormatErr
		} else if len(q.Wheres) == 1 {
			ret, err = db.Rprintf(
				"select %s from %s%s ;",
				q.ColumnNames[0],
				q.TableNames[0],
				q.Wheres[0],
			)
		} else {
			ret, err = db.Rprintf(
				"select %s from %s ;",
				q.ColumnNames[0],
				q.TableNames[0],
			)
		}
		if err != nil {
			return "", err
		}
	case RenameTableQueryType :
		if len(q.TableNames) != 2 {
			return "", NoTablesSpecifiedErr
		}

		return db.Rprintf(
			"alter table %s rename %s ;",
			q.TableNames[0],
			q.TableNames[1],
		)
	case RenameColumnQueryType :
		if len(q.ColumnNames) != 2 ||
				len(q.TableNames) != 1 {
			return "", WrongNumOfColumnsSpecifiedErr
		}

		return db.Rprintf(
			"alter table %s rename column %s to %s ;",
			q.TableNames[0],
			q.ColumnNames[0],
			q.ColumnNames[1],
		)

	case CreateTableQueryType :
		if len(q.TableSchemas) != 1 {
			return "", NoSchemaSpecifiedErr
		}

		buf, err := db.
			TableCreationStringForSchema(q.TableSchemas[0])
		if err != nil {
			return "", err
		}
		ret = Raw(buf)
		return ret, err
	case AlterColumnTypeQueryType :
		if len(q.TableNames) != 1 ||
			len(q.Columns) != 1 {
			return "",
				WrongQueryInputFormatErr
		}

		rcode, err := db.ColumnToAlterSql(
			q.Columns[0],
		)
		if err != nil {
			return "", err
		}

		buf := fmt.Sprintf(
			"alter table %s modify %s ;",
			q.TableNames[0],
			rcode,
		)

		return Raw(buf), nil
	default:
		return "", UnknownQueryTypeErr
	}

	return Raw(ret), nil
}

func (q Query)WithDB(db *DB) Query {
	q.DB = db
	return q
}

func (q Query)WithType(t QueryType) Query {
	q.Type = t
	return q
}


func (q Query)WithTableSchemas(schema ...*TableSchema) Query {
	q.TableSchemas = schema
	return q
}

func (q Query)WithWheres(where Wheres) Query {
	q.Wheres = where
	return q
}

func (q Query)WithColumnNames(
	columns ...ColumnName,
) Query {
	q.ColumnNames = columns
	return q
}

func (q Query)WithColumns(
	columns ...*Column,
) Query {
	q.Columns = columns
	return q
}

func (q Query)WithColumnTypes(
	types ...ColumnType,
) Query {
	q.ColumnTypes = types
	return q
}

func (q Query)WithTableNames(
	tables ...TableName,
) Query {
	q.TableNames = tables
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

func (q Query)RenameTable() Query {
	q.Type = RenameTableQueryType
	return q
}

func (q Query)RenameColumn() Query {
	q.Type = RenameColumnQueryType
	return q
}

func (q Query)AlterColumnType() Query {
	return q.WithType(AlterColumnTypeQueryType)
}

func (q Query)Do() (*sql.Rows, error) {
	if q.DB == nil {
		return nil, NoDBSpecifiedErr
	}
	qs, err := q.SqlRaw(q.DB)
	if err != nil {
		return nil, err
	}

	return q.DB.Query(string(qs))
}

func (v Null)SqlRaw(db *DB) (Raw, error){
	return Raw("null"), nil
}

func (v TableName)SqlRaw(db *DB) (Raw, error) {
	if v == "" {
		return "", WrongValuerFormatErr
	}
	return Raw(v), nil
}

func (v ColumnName)SqlRaw(db *DB) (Raw, error) {
	if v == "" {
		return "", WrongValuerFormatErr
	}
	return Raw(v), nil
}

func (v Raw)SqlRaw(db *DB) (Raw, error) {
	return v, nil
}

func (d Double)SqlRaw(db *DB) (Raw, error) {
	return Raw(strconv.FormatFloat(float64(d), 'f', -1, 64)),
		nil
}

func (i Int)SqlRaw(db *DB) (Raw, error) {
	return Raw(strconv.Itoa(int(i))), nil
}

func (tn TableNames)SqlRaw(db *DB) (Raw, error) {
	if len(tn) == 0 {
		return Raw(""), NoTablesSpecifiedErr
	}

	buf := make([]string, 0)
	for _, t := range tn {
		v, err := t.SqlRaw(db)
		if err != nil {
			return Raw(""), err
		}
		buf = append(buf, string(v))
	}

	ret := strings.Join(buf, ", ")
	return Raw(ret), nil
}

func (cn ColumnNames)SqlRaw(db *DB) (Raw, error) {
	if len(cn) == 0 {
		return Raw(""), NoColumnsSpecifiedErr
	}

	buf := make([]string, 0)
	for _, c := range cn {
		v, err := c.SqlRaw(db)
		if err != nil {
			return Raw(""), err
		}
		buf = append(buf, string(v))
	}

	ret := strings.Join(buf, ", ")
	return Raw(ret), nil
}

func (s String)SqlRaw(db *DB) (Raw, error) {
	ret := strings.ReplaceAll(string(s), "'", "''")
	ret = fmt.Sprintf("'%s'", ret)
	return Raw(ret), nil
}

// Return raw values separated by comma for
// column and table names and also values.
func (rvs Rawers) SqlMultiValue(db *DB) (Raw, error) {
	var ret Raw
	for i, v := range rvs {
		raw, err := v.SqlRaw(db)
		if err != nil {
			return "", err
		}

		ret += raw

		if i != len(rvs) - 1 {
			ret += ","
		}
	}

	return ret, nil
}

// Return multivalue embraced with () .
func (rvs Rawers) SqlRawTuple(db *DB) (Raw, error) {
	mval, err := rvs.SqlMultiValue(db)
	if err != nil {
		return Raw(""), err
	}

	if mval == "" {
		return "", nil
	}

	return Raw(fmt.Sprintf("(%s)", mval)), nil
}

func (db *DB)RawersEq(
	v1, v2 Rawer,
) (bool, error) {

	if v1 == nil && v2 == nil {
		return true, nil
	}

	if v1 == nil || v2 == nil {
		fmt.Println("in")
		return false, nil
	}

	raw1, err := v1.SqlRaw(db)
	if err != nil {
		return false, err
	}

	raw2, err := v2.SqlRaw(db)
	if err != nil {
		return false, err
	}
	return raw1 == raw2, nil
}

