package sqlx

import (
	"fmt"
	"errors"
	"database/sql"
	"reflect"
)

type RowValues = map[ColumnName] any

var (
	UnknownColumnType = errors.New("unknown column type")
)

func C() *Column {
	return &Column{}
}

func (c *Column) WOldName(n ColumnName) *Column {
	c.OldName = n
	return c
}

func (c *Column) WName(n ColumnName) *Column {
	c.Name = n
	return c
}

func (c *Column) WNullable() *Column {
	c.Nullable = true
	return c
}

func (c *Column) WType(t ColumnType) *Column {
	c.Type = t
	return c
}

func (c *Column) WDefault(d RawValuer) *Column {
	c.Default = d
	return c
}

func (c *Column) WKey(k Key) *Column {
	c.Key = k
	return c
}

func (c *Column) WExtra(e Code) *Column {
	c.Extra = e
	return c
}

func CT() ColumnType {
	return ColumnType{}
}

func (ct ColumnType)SqlCode(db *DB) (Code, error) {
	ret := ""
	t, ok := MysqlColumnTypeMapString[ct.VarType]
	if !ok {
		return "", UnknownColumnType
	}

	args, err := ct.Args.SqlCodeTuple(db)
	if err != nil {
		return "", err
	}

	ret = fmt.Sprintf("%s%s", t, args)

	return Code(ret), nil
}

func (ct ColumnType)Varchar(n int) ColumnType {
	ct.VarType = VarcharColumnVarType
	ct.Args = RawValuers{Int(n)}
	return ct
}

func (ct ColumnType)Int() ColumnType {
	ct.VarType = IntColumnVarType
	ct.Args = RawValuers{Int(11)}
	return ct
}

func (ct ColumnType)IntN(n int) ColumnType {
	ct.VarType = IntColumnVarType
	ct.Args = RawValuers{Int(n)}
	return ct
}

func (ct ColumnType)Nvarchar(n int) ColumnType {
	ct.VarType = NvarcharColumnVarType
	ct.Args = RawValuers{Int(n)}
	return ct
}

func (db *DB)ReadRowValues(
	rs *sql.Rows,
	ts *TableSchema,
	cnames ColumnNames,
	tsMap map[ColumnName] *Column,
	rc any,
) ([]RowValues, error) {
	null := []RowValues{}
	row := make([]any, len(cnames))
	t := reflect.TypeOf(rc)
	val := reflect.New(t)
	val = val.Elem()
	for i, v := range cnames {
		f := val.FieldByName(string(v)).Addr()
		fmt.Println("val:", val, v)
		fmt.Println("field:", f)
		c, ok := tsMap[v]
		if !ok {
			return null, ColumnDoesNotExistErr
		}
		switch c.Type.VarType {
		case VarcharColumnVarType :
			row[i] = new(sql.NullString)
		case IntColumnVarType :
			row[i] = new(sql.NullInt32)
		case TinyintColumnVarType :
			row[i] = new(sql.NullByte)
		case SmallintColumnVarType :
			row[i] = new(sql.NullInt16)
		case BigintColumnVarType :
			row[i] = new(sql.NullInt64)
		case DoubleColumnVarType :
			row[i] = new(sql.NullFloat64)
		case TimeColumnVarType :
			fallthrough
		case DatetimeColumnVarType :
			fallthrough
		case DateColumnVarType :
			fallthrough
		case TimestampColumnVarType :
			row[i] = new(sql.NullTime)
		default:
			return null, UnknownColumnTypeErr
		}
	}

	ret := []RowValues{}
	for rs.Next() {
		err := rs.Scan(row...)
		if err != nil{
			return null, err
		}
		rowMap := make(RowValues)
		for i, v := range row {
			cname := cnames[i]
			switch v.(type) {
			case *sql.NullString:
				rowMap[cname] = *(v.(*sql.NullString))
			case *sql.NullBool :
				rowMap[cname] = *(v.(*sql.NullBool))
			case *sql.NullInt32 :
				rowMap[cname] = *(v.(*sql.NullInt32))
			case *sql.NullByte:
				rowMap[cname] = *(v.(*sql.NullByte))
			case *sql.NullInt16 :
				rowMap[cname] = *(v.(*sql.NullInt16))
			case *sql.NullInt64 :
				rowMap[cname] = *(v.(*sql.NullInt64))
			case *sql.NullFloat64 :
				rowMap[cname] = *(v.(*sql.NullFloat64))
			case *sql.NullTime :
				rowMap[cname] = *(v.(*sql.NullTime))
			default:
				return null, UnknownColumnTypeErr
			}
		}
		ret = append(ret, rowMap)
	}

	return ret, nil
}

